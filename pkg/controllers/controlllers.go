package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	apisv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
	kcptenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	workloadv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/workload/v1alpha1"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"

	accessv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/access/v1alpha1"
	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/bootstrap"
	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/models"
	utilhttp "github.com/faroshq/faros-hub/pkg/util/http"
	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(accessv1alpha1.AddToScheme(scheme))
	utilruntime.Must(edgev1alpha1.AddToScheme(scheme))
	utilruntime.Must(workloadv1alpha1.AddToScheme(scheme))
	utilruntime.Must(kcptenancyv1alpha1.AddToScheme(scheme))
	utilruntime.Must(tenancyv1alpha1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
}

type Controllers interface {
	WaitForAPIReady(ctx context.Context) error
	Run(ctx context.Context) error
}

type controllerManager struct {
	config        *config.ControllerConfig
	clientFactory utilkubernetes.ClientFactory
	bootstraper   bootstrap.Bootstraper
}

func New(c *config.ControllerConfig) (Controllers, error) {
	b, err := bootstrap.New(c)
	if err != nil {
		return nil, err
	}

	cf, err := utilkubernetes.NewClientFactory(c.KCPClusterRestConfig)
	if err != nil {
		return nil, err
	}

	return &controllerManager{
		config:        c,
		clientFactory: cf,
		bootstraper:   b,
	}, nil
}

func (c *controllerManager) WaitForAPIReady(ctx context.Context) error {
	// Wait for API server to report healthy
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		h := utilhttp.GetInsecureClient()
		res, err := h.Get(c.config.KCPClusterRestConfig.Host + "/healthz")
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

func (c *controllerManager) Run(ctx context.Context) error {
	// bootstrap will set missing ctrlRestConfig and deploy kcp wide resources
	ctxT, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	plugins, err := c.bootstrap(ctxT)
	if err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}

	eg := errgroup.Group{}

	eg.Go(func() error {
		return c.runEdge(ctx, plugins)
	})
	eg.Go(func() error {
		return c.runSystem(ctx)
	})

	return eg.Wait()
}

func (c *controllerManager) bootstrap(ctx context.Context) (models.PluginsList, error) {
	// create controllers workspace
	for _, w := range []string{
		c.config.TenantsWorkspacePrefix,
		c.config.ControllersTenantWorkspace,
		c.config.ControllersWorkspace,
		c.config.ControllersPluginsWorkspace,
	} {
		if err := c.bootstraper.CreateWorkspace(ctx, w); err != nil {
			return nil, err
		}
	}
	// create assets for controller workspace being able to access all "workspaces"
	// and implement their requests
	if err := c.bootstraper.DeployKustomizeAssetsCRD(ctx, c.config.ControllersWorkspace); err != nil {
		return nil, err
	}
	if err := c.bootstraper.DeployKustomizeAssetsKCP(ctx, c.config.ControllersWorkspace); err != nil {
		return nil, err
	}

	// create assets for controller tenant workspace being able to access use apis
	if err := c.bootstraper.BootstrapServiceTenantAssets(ctx, c.config.ControllersTenantWorkspace); err != nil {
		return nil, err
	}

	// load plugins and create assets for each tenant
	plugins, err := c.bootstraper.LoadPlugins(ctx, c.config.ControllersPluginsWorkspace)
	if err != nil {
		return nil, err
	}

	return plugins, nil
}

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
