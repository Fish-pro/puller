package app

import (
	"context"
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/puller.io/puller/cmd/puller/app/options"
)

func NewControllerManagerCommand(ctx context.Context) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "puller-controller-manager",
		Short: "puller-controller-manager is a controller manager for pull private image",
		RunE: func(cmd *cobra.Command, args []string) error {
			if errs := opts.Validate(); len(errs) != 0 {
				return errs.ToAggregate()
			}
			return Run(ctx, opts)
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
	genericFlagSet.Lookup("kubeconfig").Usage = "Path to kubeconfig file."
	opts.AddFlags(genericFlagSet, []string{})

	return cmd
}

func Run(ctx context.Context, opts *options.Options) error {
	return nil
}
