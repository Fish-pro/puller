package options

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/component-base/logs"
)

type Options struct {
	// Controllers contains all controller names.
	Controllers []string
	// metricsAddr define the metrics addr
	MetricsAddr string
	// enableLeaderElection define enable leader election
	EnableLeaderElection bool
	// probeAddr define the probe address
	ProbeAddr string
	// KubeAPIQPS is the QPS to use while talking with karmada-apiserver.
	KubeAPIQPS float32
	// KubeAPIBurst is the burst to allow while talking with karmada-apiserver.
	KubeAPIBurst int
	// ConcurrentPullerSyncs is the number of puller objects that are
	// allowed to sync concurrently.
	ConcurrentPullerSyncs int

	Logs *logs.Options

	Master     string
	Kubeconfig string
}

func NewOptions() *Options {
	return &Options{
		Logs: logs.NewOptions(),
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet, allControllers []string) {
	if o == nil {
		return
	}
	fs.StringSliceVar(&o.Controllers, "controllers", []string{"*"}, fmt.Sprintf(
		"A list of controllers to enable. '*' enables all on-by-default controllers, 'foo' enables the controller named 'foo', '-foo' disables the controller named 'foo'. All controllers: %s.",
		strings.Join(allControllers, ", "),
	))
	fs.StringVar(&o.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	fs.StringVar(&o.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	fs.BoolVar(&o.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	fs.Float32Var(&o.KubeAPIQPS, "kube-api-qps", 40.0, "QPS to use while talking with karmada-apiserver. Doesn't cover events and node heartbeat apis which rate limiting is controlled by a different set of flags.")
	fs.IntVar(&o.KubeAPIBurst, "kube-api-burst", 60, "Burst to use while talking with karmada-apiserver. Doesn't cover events and node heartbeat apis which rate limiting is controlled by a different set of flags.")
	fs.IntVar(&o.ConcurrentPullerSyncs, "concurrent-puller-syncs", 5, "The number of Puller that are allowed to sync concurrently.")
	fs.StringVar(&o.Master, "master", o.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig).")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")

}
