package main

import (
	"os"

	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register"
	controllerruntime "sigs.k8s.io/controller-runtime"
	_ "sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/puller.io/puller/cmd/puller/app"
)

func main() {
	ctx := controllerruntime.SetupSignalHandler()
	cmd := app.NewControllerManagerCommand(ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
