package main

import (
	"context"
	"flag"
	"strings"

	devproxyclient "github.com/faroshq/faros-hub/pkg/dev/proxy-client"
	devproxyserver "github.com/faroshq/faros-hub/pkg/dev/proxy-server"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	certFile      = flag.String("certFile", "dev/proxy.crt", "file containing server certificate")
	keyFile       = flag.String("keyFile", "dev/proxy.key", "file containing server key")
	serverAddress = flag.String("serverAddress", "0.0.0.0:8443", "Server address")

	clientCertFile      = flag.String("clientCertFile", "dev/proxy-client.crt", "file containing client certificate")
	clientCertKeyFile   = flag.String("clientCertKeyFile", "dev/proxy-client.key", "file containing client key")
	clientUpstreamURL   = flag.String("clientUpstreamUrl", "https://localhost:8443", "Server external address")
	clientDownstreamURL = flag.String("clientDownstreamUrl", "http://localhost:8080", "Client forward address")
	clientID            = flag.String("clientID", "faros-dev", "Client ID")
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
	case "client":
		return runClient(ctx)
	default:
		return runServer(ctx)
	}
}

func runServer(ctx context.Context) error {
	server, err := devproxyserver.New(*serverAddress, *certFile, *keyFile)
	if err != nil {
		return err
	}
	return server.Run(ctx)
}

func runClient(ctx context.Context) error {
	client, err := devproxyclient.New(*clientUpstreamURL, *clientDownstreamURL, *clientCertFile, *clientCertKeyFile, *clientID)
	if err != nil {
		return err
	}
	go client.Run(ctx)
	<-ctx.Done()
	return nil
}
