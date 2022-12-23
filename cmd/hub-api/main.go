package main

import (
	"context"
	"flag"
	"os"

	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/server"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
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
	c, err := config.LoadAPI()
	if err != nil {
		return err
	}

	server, err := server.New(ctx, c)
	if err != nil {
		return err
	}
	return server.Run(ctx)
}
