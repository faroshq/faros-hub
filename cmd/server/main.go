package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/kcp-dev/kcp/pkg/embeddedetcd"
	kcpfeatures "github.com/kcp-dev/kcp/pkg/features"
	"github.com/kcp-dev/kcp/pkg/server"
	"github.com/kcp-dev/kcp/pkg/server/options"
	"github.com/mjudeikis/kcp-example/pkg/bootstrap"
	"github.com/mjudeikis/kcp-example/pkg/config"
	fconfig "github.com/mjudeikis/kcp-example/pkg/config"
	utilhttp "github.com/mjudeikis/kcp-example/pkg/util/http"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	genericapiserver "k8s.io/apiserver/pkg/server"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	kconfig "k8s.io/component-base/config"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var skipKCP bool

func init() {
	flag.BoolVar(&skipKCP, "skip-kcp", false, "skip KCP server")
}

func main() {
	ctx := genericapiserver.SetupSignalContext()
	flag.Parse()
	if err := run(ctx); err != nil {
		fmt.Printf("error starting server: %v", err)
		os.Exit(1)
	}
}

func runKCP(ctx context.Context, c *config.Config) error {
	// Create a new health instance
	h := health.New()
	defer h.Stop()

	logger := klog.FromContext(ctx).WithValues("component", "kcp")
	ctx = klog.NewContext(ctx, logger)

	serverOptions := options.NewOptions(c.Server.StateDir)
	serverOptions.GenericControlPlane.Logs.Config.Verbosity = kconfig.VerbosityLevel(5)

	// set tunnels true
	runtime.Must(utilfeature.DefaultMutableFeatureGate.Set(fmt.Sprintf("%s=true", kcpfeatures.SyncerTunnel)))

	completed, err := serverOptions.Complete()
	if err != nil {
		return err
	}

	if errs := completed.Validate(); len(errs) > 0 {
		return errors.NewAggregate(errs)
	}

	serverConfig, err := server.NewConfig(completed)
	if err != nil {
		return err
	}

	completedConfig, err := serverConfig.Complete()
	if err != nil {
		return err
	}

	klog.Infof("Batteries included: %s", strings.Join(completed.Extra.BatteriesIncluded, ","))

	if completedConfig.EmbeddedEtcd.Config != nil {
		if err := embeddedetcd.NewServer(completedConfig.EmbeddedEtcd).Run(ctx); err != nil {
			return err
		}
	}

	s, err := server.NewServer(completedConfig)
	if err != nil {
		return err
	}

	// Run KCP service
	go s.Run(ctx)
	return nil
}

func run(ctx context.Context) error {
	c, err := fconfig.Load(true)
	if err != nil {
		return err
	}

	if !skipKCP {
		if err := runKCP(ctx, c); err != nil {
			return err
		}
	}

	// Wait for API server to report healthy
	done := false
	for !done {
		h := utilhttp.GetInsecureClient()
		res, err := h.Get("https://localhost:6443/healthz")
		switch {
		case err != nil:
			klog.Infof("Waiting for API server to report healthy: %v", err)
		case res.StatusCode != http.StatusOK:
			klog.Infof("Waiting for API server to report healthy: %v", res.Status)
		case res.StatusCode == http.StatusOK:
			klog.Infof("API server is healthy")
			done = true
		}
		time.Sleep(1 * time.Second)
	}

	klog.Infof("KCP is ready. Bootstrapping store...")
	os.Setenv("KUBECONFIG", filepath.Join(c.Server.StateDir, "admin.kubeconfig"))
	restConfig := ctrl.GetConfigOrDie()

	b, err := bootstrap.New(ctx, c, restConfig)
	if err != nil {
		return err
	}

	err = b.BootstrapOrganization(ctx)
	if err != nil {
		return err
	}

	err = b.BootstrapCompute(ctx)
	if err != nil {
		return err
	}

	err = b.BootstrapServices(ctx)
	if err != nil {
		return err
	}

	err = b.BootstrapUsers(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}
