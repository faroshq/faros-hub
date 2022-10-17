package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/faroshq/faros-hub/pkg/controllers/edge/registration"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/kcp"
)

// edge controller is running in edge api virtual workspace context

func (c *controllers) runEdge(ctx context.Context) error {
	restConfig, err := c.clientFactory.GetWorkspaceRestConfig(ctx, c.config.ControllersFarosEdgeAPIExportName)
	if err != nil {
		return err
	}

	var edgeRest *rest.Config
	// bootstrap rest config for controllers
	if kcpAPIsGroupPresent(restConfig) {
		if err := wait.PollImmediateInfinite(time.Second*5, func() (bool, error) {
			klog.Info("looking up virtual workspace URL - edge")
			edgeRest, err = restConfigForAPIExport(ctx, restConfig, c.config.ControllersFarosEdgeAPIExportName)
			if err != nil {
				return false, nil
			}
			return true, nil
		}); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("kcp APIs group not present in cluster. We don't support non kcp clusters yet")
	}

	options := ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      ":8080",
		Port:                    9443,
		HealthProbeBindAddress:  ":8081",
		LeaderElection:          false,
		LeaderElectionID:        "edge.faros.sh",
		LeaderElectionNamespace: "default",
		LeaderElectionConfig:    edgeRest,
	}

	mgr, err := kcp.NewClusterAwareManager(edgeRest, options)
	if err != nil {
		klog.Error(err, "unable to start manager")
		return err
	}

	if err = (&registration.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		klog.Error(err, "unable to create controller", "registration.edge.faros.sh")
		return err
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		klog.Error(err, "unable to set up health check")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		klog.Error(err, "unable to set up ready check")
		return err
	}

	klog.Info("starting edge manager")

	return mgr.Start(ctx)

}
