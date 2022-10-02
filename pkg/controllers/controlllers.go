package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	apisv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/kcp"

	accessv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/access/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/bootstrap"
	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/controllers/access"
	utilhttp "github.com/faroshq/faros-hub/pkg/util/http"
	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(accessv1alpha1.AddToScheme(scheme))
}

type Controllers interface {
	WaitForAPIReady(ctx context.Context) error
	Run(ctx context.Context) error
}

type controllers struct {
	config        *config.Config
	clientFactory utilkubernetes.ClientFactory

	// ctrlRestConfig is workspaces rest config for controllers workspace
	// for operators to be operating in
	ctrlRestConfig *rest.Config
	bootstraper    bootstrap.Bootstraper
}

func New(c *config.Config, r *rest.Config) (Controllers, error) {
	b, err := bootstrap.New(c, r)
	if err != nil {
		return nil, err
	}

	cf, err := utilkubernetes.NewClientFactory(r)
	if err != nil {
		return nil, err
	}

	return &controllers{
		config:        c,
		clientFactory: cf,
		bootstraper:   b,
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
		res, err := h.Get("https://localhost:6443/healthz")
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
	// bootstrap will set missing ctrlRestConfig and deploy kcp wide resources
	if err := c.bootstrap(ctx); err != nil {
		return err
	}
	time.Sleep(time.Second * 5)

	_, rest, err := c.clientFactory.GetWorkspaceClient(ctx, c.config.ControllersWorkspace)
	if err != nil {
		return err
	}

	options := ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      ":8080",
		Port:                    9443,
		HealthProbeBindAddress:  ":8081",
		LeaderElection:          false,
		LeaderElectionID:        "controllers.faros.sh",
		LeaderElectionNamespace: "default",
		LeaderElectionConfig:    rest,
	}

	mgr, err := kcp.NewClusterAwareManager(c.ctrlRestConfig, options)
	if err != nil {
		klog.Error(err, "unable to start manager")
		return err
	}

	if err = (&access.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Config: c.config,
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

func (c *controllers) bootstrap(ctx context.Context) error {
	// create controllers workspace
	if err := c.bootstraper.CreateWorkspace(ctx, c.config.ControllersWorkspace); err != nil {
		return err
	}
	// create assets for controller workspace being able to access all "workspaces"
	// and implement their requests
	if err := c.bootstraper.DeployKustomizeAssets(ctx, c.config.ControllersWorkspace, "./config/crd"); err != nil {
		return err
	}
	if err := c.bootstraper.DeployKustomizeAssets(ctx, c.config.ControllersWorkspace, "./config/kcp"); err != nil {
		return err
	}

	// controller would need to be running in workspace context, so we need to get rest.Config
	// for that workspace
	_, rest, err := c.clientFactory.GetWorkspaceClient(ctx, c.config.ControllersWorkspace)
	if err != nil {
		return err
	}

	// bootstrap rest config for controllers
	if kcpAPIsGroupPresent(rest) {
		klog.Info("Looking up virtual workspace URL")
		rest, err := restConfigForAPIExport(ctx, rest, c.config.ControllersAPIExport)
		if err != nil {
			return err
		}
		c.ctrlRestConfig = rest
	} else {
		return fmt.Errorf("kcp APIs group not present in cluster. We don't support non kcp clusters yet")
	}

	return nil
}

// +kubebuilder:rbac:groups="apis.kcp.dev",resources=apiexports,verbs=get;list;watch

// restConfigForAPIExport returns a *rest.Config properly configured to communicate with the endpoint for the
// APIExport's virtual workspace.
func restConfigForAPIExport(ctx context.Context, cfg *rest.Config, apiExportName string) (*rest.Config, error) {
	scheme := runtime.NewScheme()
	if err := apisv1alpha1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("error adding apis.kcp.dev/v1alpha1 to scheme: %w", err)
	}

	apiExportClient, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("error creating APIExport client: %w", err)
	}

	var apiExport apisv1alpha1.APIExport

	if apiExportName != "" {
		if err := apiExportClient.Get(ctx, types.NamespacedName{Name: apiExportName}, &apiExport); err != nil {
			return nil, fmt.Errorf("error getting APIExport %q: %w", apiExportName, err)
		}
	} else {
		klog.Infof("api-export-name is empty - listing")
		exports := &apisv1alpha1.APIExportList{}
		if err := apiExportClient.List(ctx, exports); err != nil {
			return nil, fmt.Errorf("error listing APIExports: %w", err)
		}
		if len(exports.Items) == 0 {
			return nil, fmt.Errorf("no APIExport found")
		}
		if len(exports.Items) > 1 {
			return nil, fmt.Errorf("more than one APIExport found")
		}
		apiExport = exports.Items[0]
	}

	if len(apiExport.Status.VirtualWorkspaces) < 1 {
		return nil, fmt.Errorf("APIExport %q status.virtualWorkspaces is empty", apiExportName)
	}

	cfg = rest.CopyConfig(cfg)
	cfg.Host = apiExport.Status.VirtualWorkspaces[0].URL

	return cfg, nil
}

func kcpAPIsGroupPresent(restConfig *rest.Config) bool {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		klog.Error(err, "failed to create discovery client")
		os.Exit(1)
	}
	apiGroupList, err := discoveryClient.ServerGroups()
	if err != nil {
		klog.Error(err, "failed to get server groups")
		os.Exit(1)
	}

	for _, group := range apiGroupList.Groups {
		if group.Name == apisv1alpha1.SchemeGroupVersion.Group {
			for _, version := range group.Versions {
				if version.Version == apisv1alpha1.SchemeGroupVersion.Version {
					return true
				}
			}
		}
	}
	return false
}
