package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	devproxyclient "github.com/faroshq/faros-hub/pkg/dev/proxy-client"
	devproxyserver "github.com/faroshq/faros-hub/pkg/dev/proxy-server"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	certFile            = flag.String("certFile", "dev/proxy.crt", "file containing server certificate")
	keyFile             = flag.String("keyFile", "dev/proxy.key", "file containing server key")
	clientCertFile      = flag.String("clientCertFile", "dev/proxy-client.crt", "file containing client certificate")
	clientCertKeyFile   = flag.String("clientCertKeyFile", "dev/proxy-client.key", "file containing client key")
	clientUpstreamURL   = flag.String("client-upstream-url", "https://localhost:8443", "Server external address")
	clientDownstreamURL = flag.String("client-downstream-url", "http://localhost:9090", "Client forward address")
	serverAddress       = flag.String("server-address", "localhost:8443", "Server address")
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
		panic(err)
	}
}

func run(ctx context.Context) error {
	switch strings.ToLower(flag.Arg(0)) {
	case "server":
		return runServer(ctx)
	case "client":
		return runClient(ctx)
	default:
		return fmt.Errorf("unknown mode %q", flag.Arg(0))
	}
}

func runServer(ctx context.Context) error {
	server, err := devproxyserver.New(*serverAddress, *certFile, *keyFile, *clientCertFile)
	if err != nil {
		return err
	}
	return server.Run(ctx)
}

func runClient(ctx context.Context) error {
	client, err := devproxyclient.New(*clientUpstreamURL, *clientDownstreamURL, *clientCertFile, *clientCertKeyFile, *certFile)
	if err != nil {
		return err
	}
	go client.Run(ctx)
	<-ctx.Done()
	return nil
}
