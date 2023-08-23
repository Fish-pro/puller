package app

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlconfig "sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/puller-io/puller/cmd/puller/app/options"
	pullerv1alpha1 "github.com/puller-io/puller/pkg/apis/puller/v1alpha1"
	"github.com/puller-io/puller/pkg/controller/puller"
	"github.com/puller-io/puller/pkg/scheme"
	"github.com/puller-io/puller/pkg/version"
)

func NewControllerManagerCommand(ctx context.Context) *cobra.Command {
	o := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "puller",
		Short: "puller is a controller for pull private image",
		RunE: func(cmd *cobra.Command, args []string) error {
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
	genericFlagSet.Lookup("kubeconfig").Usage = "Paths to a kubeconfig. Only required if out-of-cluster."
	o.AddFlags(genericFlagSet)

	cmd.Flags().AddFlagSet(genericFlagSet)
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
		Logger:                     logger,
		Scheme:                     scheme.Scheme,
		SyncPeriod:                 &opts.ResyncPeriod.Duration,
		LeaderElection:             opts.LeaderElection.LeaderElect,
		LeaderElectionID:           opts.LeaderElection.ResourceName,
		LeaderElectionNamespace:    opts.LeaderElection.ResourceNamespace,
		LeaseDuration:              &opts.LeaderElection.LeaseDuration.Duration,
		RenewDeadline:              &opts.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:                &opts.LeaderElection.RetryPeriod.Duration,
		LeaderElectionResourceLock: opts.LeaderElection.ResourceLock,
		MetricsBindAddress:         opts.MetricsAddr,
		HealthProbeBindAddress:     opts.ProbeAddr,
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

	if err = (&puller.Controller{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		KubeClient:    kubernetes.NewForConfigOrDie(mgr.GetConfig()),
		EventRecorder: mgr.GetEventRecorderFor(puller.ControllerName),
	}).SetupWithManager(mgr); err != nil {
		klog.Error(err, "unable to create controller", "controller", "Puller")
		return fmt.Errorf("create puller controller failed, error: %v", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		klog.Errorf("unable to set up health check: %w", err)
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		klog.Error("unable to set up ready check: %w", err)
		return err
	}

	// blocks until the context is done.
	if err := mgr.Start(ctx); err != nil {
		klog.Errorf("controller manager exits unexpectedly: %v", err)
		return err
	}
	return nil
}
