package app

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	cliflag "k8s.io/component-base/cli/flag"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlconfig "sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/puller-io/puller/cmd/puller/app/options"
	pullerv1alpha1 "github.com/puller-io/puller/pkg/apis/puller/v1alpha1"
	clientbuilder "github.com/puller-io/puller/pkg/builder"
	"github.com/puller-io/puller/pkg/controller/puller"
	"github.com/puller-io/puller/pkg/util/gclient"
	"github.com/puller-io/puller/pkg/version"
)

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
}

func NewControllerManagerCommand(ctx context.Context) *cobra.Command {
	o := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "puller-controller-manager",
		Short: "puller-controller-manager is a controller manager for pull private image",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Activate logging as soon as possible, after that
			// show flags with the final logging configuration.
			if err := logsapi.ValidateAndApply(o.Logs, utilfeature.DefaultFeatureGate); err != nil {
				return err
			}
			cliflag.PrintFlags(cmd.Flags())

			if errs := o.Validate(); len(errs) != 0 {
				return errs.ToAggregate()
			}
			return Run(ctx, o)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fss := cliflag.NamedFlagSets{}
	genericFlagSet := fss.FlagSet("generic")
	genericFlagSet.AddGoFlagSet(flag.CommandLine)
	o.AddFlags(genericFlagSet, controllers.ControllerNames())

	logsapi.AddFlags(o.Logs, fss.FlagSet("logs"))
	return cmd
}

func Run(ctx context.Context, opts *options.Options) error {
	logger := klog.FromContext(ctx)

	// To help debugging, immediately log version
	logger.Info("Starting", "version", version.Get())
	logger.Info("Golang settings", "GOGC", os.Getenv("GOGC"), "GOMAXPROCS", os.Getenv("GOMAXPROCS"), "GOTRACEBACK", os.Getenv("GOTRACEBACK"))

	config, err := ctrl.GetConfig()
	if err != nil {
		logger.Error(err, "failed to get config")
		return err
	}
	config.QPS, config.Burst = opts.KubeAPIQPS, opts.KubeAPIBurst
	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Logger:                 logger,
		Scheme:                 gclient.NewSchema(),
		MetricsBindAddress:     opts.MetricsAddr,
		HealthProbeBindAddress: opts.ProbeAddr,
		LeaderElection:         opts.EnableLeaderElection,
		LeaderElectionID:       "688adee7.puller.io",
		BaseContext: func() context.Context {
			return ctx
		},
		Controller: ctrlconfig.Controller{
			GroupKindConcurrency: map[string]int{
				pullerv1alpha1.SchemeGroupVersion.WithKind("Puller").GroupKind().String(): opts.ConcurrentPullerSyncs,
			},
		},
	})
	if err != nil {
		klog.Errorf("Failed to build controller manager: %w", err)
		return err
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		klog.Errorf("unable to set up health check: %w", err)
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		klog.Error("unable to set up ready check: %w", err)
		return err
	}

	setupControllers(mgr, opts, ctx.Done())

	// blocks until the context is done.
	if err := mgr.Start(ctx); err != nil {
		klog.Errorf("controller manager exits unexpectedly: %v", err)
		return err
	}
	return nil
}

var (
	// controllersDisabledByDefault is the set of controllers which is disabled by default
	controllersDisabledByDefault = sets.Set[string]{}
	controllers                  = make(Initializers)
)

// ControllerContext defines the context object for controller.
type ControllerContext struct {
	Mgr              ctrl.Manager
	Opts             options.Options
	StopChan         <-chan struct{}
	KubeClientSet    kubernetes.Interface
	DynamicClientSet dynamic.Interface
	ClientBuilder    clientbuilder.ControllerClientBuilder
}

func init() {
	controllers["puller"] = StartPullerController
}

// InitFunc is used to launch a particular controller.
// Any error returned will cause the controller process to `Fatal`
// The bool indicates whether the controller was enabled.
type InitFunc func(ctx ControllerContext) (enabled bool, err error)

// Initializers is a public map of named controller groups
type Initializers map[string]InitFunc

// ControllerNames returns all known controller names
func (i Initializers) ControllerNames() []string {
	return sets.StringKeySet(i).List()
}

// IsControllerEnabled check if a specified controller enabled or not.
func (c ControllerContext) IsControllerEnabled(name string, disabledByDefaultControllers sets.Set[string]) bool {
	hasStar := false
	for _, ctrl := range c.Opts.Controllers {
		if ctrl == name {
			return true
		}
		if ctrl == "-"+name {
			return false
		}
		if ctrl == "*" {
			hasStar = true
		}
	}
	if !hasStar {
		return false
	}

	return !disabledByDefaultControllers.Has(name)
}

// StartControllers starts a set of controllers with a specified ControllerContext
func (i Initializers) StartControllers(ctx ControllerContext, controllersDisabledByDefault sets.Set[string]) error {
	for controllerName, initFn := range i {
		if !ctx.IsControllerEnabled(controllerName, controllersDisabledByDefault) {
			klog.Warningf("%q is disabled", controllerName)
			continue
		}
		klog.V(1).Infof("Starting %q", controllerName)
		started, err := initFn(ctx)
		if err != nil {
			klog.Errorf("Error starting %q", controllerName)
			return err
		}
		if !started {
			klog.Warningf("Skipping %q", controllerName)
			continue
		}
		klog.Infof("Started %q", controllerName)
	}
	return nil
}

// setupControllers initialize controllers and setup one by one.
func setupControllers(mgr ctrl.Manager, opts *options.Options, stopChan <-chan struct{}) {
	restConfig := mgr.GetConfig()
	kubeClientSet := kubernetes.NewForConfigOrDie(restConfig)
	dynamicClientSet := dynamic.NewForConfigOrDie(restConfig)
	clientBuilder := clientbuilder.NewPullerControllerClientBuilder(restConfig)

	controllerContext := ControllerContext{
		Mgr:              mgr,
		Opts:             *opts,
		StopChan:         stopChan,
		KubeClientSet:    kubeClientSet,
		DynamicClientSet: dynamicClientSet,
		ClientBuilder:    clientBuilder,
	}

	if err := controllers.StartControllers(controllerContext, controllersDisabledByDefault); err != nil {
		klog.Fatalf("error starting controllers: %v", err)
	}
}

func StartPullerController(ctx ControllerContext) (enabled bool, err error) {
	pullerController := &puller.Controller{
		Client:        ctx.Mgr.GetClient(),
		KubeClient:    ctx.KubeClientSet,
		EventRecorder: ctx.Mgr.GetEventRecorderFor(puller.ControllerName),
	}
	if err := pullerController.SetupWithManager(ctx.Mgr); err != nil {
		return false, err
	}
	return true, nil
}
