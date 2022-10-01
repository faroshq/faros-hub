package main

import (
	"context"
	"flag"
	"os"

	"k8s.io/klog/v2"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/controllers"
)

func main() {
	opts := zap.Options{
		Development: true,
	}

	opts.BindFlags(flag.CommandLine)
	klog.InitFlags(flag.CommandLine)

	flag.Parse()
	flag.Lookup("v").Value.Set("6")

	ctx := ctrl.SetupSignalHandler()

	err := run(ctx)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	restConfig := ctrl.GetConfigOrDie()

	c, err := config.Load()
	if err != nil {
		return err
	}

	controllers, err := controllers.New(c, restConfig)
	if err != nil {
		return err
	}

	err = controllers.WaitForAPIReady(ctx)
	if err != nil {
		return err
	}

	return controllers.Run(ctx)

}
