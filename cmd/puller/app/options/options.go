package options

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	componentbaseconfig "k8s.io/component-base/config"
)

type Options struct {
	// Controllers contains all controller names.
	Controllers    []string
	LeaderElection componentbaseconfig.LeaderElectionConfiguration
	// BindAddress is the IP address on which to listen for the --secure-port port.
	BindAddress string
	// SecurePort is the port that the the server serves at.
	// Note: We hope support https in the future once controller-runtime provides the functionality.
	SecurePort int
}

func NewOptions() *Options {
	return &Options{
		LeaderElection: componentbaseconfig.LeaderElectionConfiguration{
			LeaderElect:  true,
			ResourceLock: resourcelock.LeasesResourceLock,
		},
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
}
