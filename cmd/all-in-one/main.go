package main

import (
	"context"
	"flag"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/controllers"
	"github.com/faroshq/faros-hub/pkg/server"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	opts := zap.Options{
		Development: true,
	}

	opts.BindFlags(flag.CommandLine)
	klog.InitFlags(flag.CommandLine)

	flag.Parse()
	flag.Lookup("v").Value.Set("6")

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ctx := ctrl.SetupSignalHandler()
	ctx = klog.NewContext(ctx, klog.NewKlogr())

	err := run(ctx)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cc, err := config.LoadController()
	if err != nil {
		return err
	}

	controllers, err := controllers.New(cc)
	if err != nil {
		return err
	}

	err = controllers.WaitForAPIReady(ctx)
	if err != nil {
		return err
	}

	ca, err := config.LoadAPI()
	if err != nil {
		return err
	}

	server, err := server.New(ctx, ca)
	if err != nil {
		return err
	}
	go controllers.Run(ctx)
	return server.Run(ctx)
}
