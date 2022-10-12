package controllers

import (
	"context"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/edge/controllers/agent"
	utilhttp "github.com/faroshq/faros-hub/pkg/util/http"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(edgev1alpha1.AddToScheme(scheme))
}

type Controllers interface {
	WaitForAPIReady(ctx context.Context) error
	Run(ctx context.Context) error
}

type controllers struct {
	config *config.Config

	rest *rest.Config
}

func New(c *config.Config) (Controllers, error) {
	restConfig := ctrl.GetConfigOrDie()

	return &controllers{
		config: c,
		rest:   restConfig,
	}, nil
}

func (c *controllers) WaitForAPIReady(ctx context.Context) error {
	// Wait for API server to report healthy
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		h := utilhttp.GetInsecureClient()
		res, err := h.Get(c.rest.Host + "/healthz")
		switch {
		case err != nil:
			klog.Infof("Waiting for API server to report healthy: %v", err)
		case res.StatusCode != http.StatusOK:
			klog.Infof("Waiting for API server to report healthy: %v", res.Status)
		case res.StatusCode == http.StatusOK:
			klog.Infof("API server is healthy")
			return nil
		}

		select {
		case <-ctx.Done():
			klog.Infof("stopped waiting for API server to report healthy: %v", ctx.Err())
			return nil
		case <-ticker.C:
		}
	}
}

func (c *controllers) Run(ctx context.Context) error {
	options := ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      ":8080",
		Port:                    9443,
		HealthProbeBindAddress:  ":8081",
		LeaderElection:          false,
		LeaderElectionID:        "edge.faros.sh",
		LeaderElectionNamespace: "default",
		LeaderElectionConfig:    c.rest,
	}

	mgr, err := ctrl.NewManager(c.rest, options)
	if err != nil {
		return err
	}

	coreClients, err := kubernetes.NewClusterForConfig(c.config.RootRestConfig)
	if err != nil {
		return err
	}

	if err = (&agent.Reconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		Config:      c.config,
		CoreClients: coreClients,
	}).SetupWithManager(mgr); err != nil {
		klog.Error(err, "unable to create controller", "controller")
		return err
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		klog.Error(err, "unable to set up health check")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		klog.Error(err, "unable to set up ready check")
		return err
	}

	klog.Info("starting manager")

	return mgr.Start(ctx)
}
