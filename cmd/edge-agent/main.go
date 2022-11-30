package main

import (
	"context"
	"flag"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/edge/controllers"
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
	ctx = klog.NewContext(ctx, klog.NewKlogr())

	err := run(ctx)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	c, err := config.LoadAgent()
	if err != nil {
		return err
	}

	controllers, err := controllers.New(c)
	if err != nil {
		return err
	}

	err = controllers.WaitForAPIReady(ctx)
	if err != nil {
		return err
	}

	return controllers.Run(ctx)

}
