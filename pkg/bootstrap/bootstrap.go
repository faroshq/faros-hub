package bootstrap

import (
	"context"
	"net/url"
	"path/filepath"
	"strings"

	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/mjudeikis/kcp-example/pkg/config"
	utilkubernetes "github.com/mjudeikis/kcp-example/pkg/util/kubernetes"
	"k8s.io/client-go/rest"
)

type Bootstraper interface {
	BootstrapOrganization(ctx context.Context) error
	BootstrapCompute(ctx context.Context) error
	BootstrapServices(ctx context.Context) error
	BootstrapUsers(ctx context.Context) error
}

type bootstrap struct {
	rest   *rest.Config
	config *config.Config

	kcpClient kcpclient.ClusterInterface
}

func New(ctx context.Context, config *config.Config, rest *rest.Config) (*bootstrap, error) {
	kcpClient, err := newKCPClusterClient(rest)
	if err != nil {
		return nil, err
	}

	b := &bootstrap{
		rest:      rest,
		config:    config,
		kcpClient: kcpClient,
	}

	return b, nil
}

func (b *bootstrap) BootstrapOrganization(ctx context.Context) error {
	for _, workspace := range b.config.Server.WorkspacesList {
		err := b.createNamedWorkspace(ctx, workspace)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *bootstrap) BootstrapCompute(ctx context.Context) error {
	// Read paths of compute kubeconfigs for services and shared
	// create syncers and deploy into clusters
	// Add Location and Placement for individual services
	if len(b.config.Server.ComputeServicesKubeconfigs) == 0 && len(b.config.Server.ComputeSharedKubeconfigs) == 0 {
		return nil
	}

	labelsService := map[string]string{
		"role": "service",
	}

	for _, kubeconfig := range b.config.Server.ComputeServicesKubeconfigs {
		name := strings.Split(filepath.Base(kubeconfig), ".")[0]
		rest, err := utilkubernetes.GetRestConfigFromURL(kubeconfig)
		if err != nil {
			return err
		}
		err = b.bootstrapSyncTargets(ctx, name, b.config.Server.ComputeServiceWorkspace, rest, labelsService)
		if err != nil {
			return err
		}
	}

	// Locations defines locations as clusters
	err := b.bootstrapLocations(ctx, "services", b.config.Server.ComputeServiceWorkspace, labelsService)
	if err != nil {
		return err
	}

	// Placements targets locations via selectors so it know where to put workloads in
	err = b.bootstrapPlacements(ctx, "root:"+b.config.Server.ComputeServiceWorkspace, b.config.Server.ControllerServicesWorkspace, "services", labelsService)
	if err != nil {
		return err
	}

	// As our workspace for compute is different from one where services are running,
	// we need to expose kubernetes resources into those.
	// By default we get `kubernetes` APIExport created, we ned to bind to it
	err = b.bootstrapBinding(ctx, "root:"+b.config.Server.ComputeServiceWorkspace, "kubernetes", b.config.Server.ControllerServicesWorkspace, "kubernetes")
	if err != nil {
		return err
	}

	labelsShared := map[string]string{
		"role": "shared",
	}

	for _, kubeconfig := range b.config.Server.ComputeSharedKubeconfigs {
		name := strings.Split(filepath.Base(kubeconfig), ".")[0]
		rest, err := utilkubernetes.GetRestConfigFromURL(kubeconfig)
		if err != nil {
			return err
		}
		err = b.bootstrapSyncTargets(ctx, name, b.config.Server.ComputeSharedWorkspace, rest, labelsShared)
		if err != nil {
			return err
		}
	}

	err = b.bootstrapLocations(ctx, "shared", b.config.Server.ComputeSharedWorkspace, labelsService)
	if err != nil {
		return err
	}

	// As our workspace for compute is different from one where services are running,
	// we need to expose kubernetes resources into those.
	// By default we get `kubernetes` APIExport created, we ned to bind to it
	err = b.bootstrapBinding(ctx, "root:"+b.config.Server.ControllerServicesWorkspace, "faros.sh", "users:user1", "faros.sh")
	if err != nil {
		return err
	}

	err = b.bootstrapBinding(ctx, "root:"+b.config.Server.ComputeSharedWorkspace, "kubernetes", "users:user1", "kubernetes")
	if err != nil {
		return err
	}

	return nil
}

func (b *bootstrap) BootstrapServices(ctx context.Context) error {

	toDeploy := []string{
		"./config/kcp",
		"./config/crd",
		"./config/manager",
		"./config/rbac",
	}

	for _, path := range toDeploy {
		err := b.deployComponents(ctx, b.config.Server.ControllerServicesWorkspace, path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *bootstrap) BootstrapUsers(ctx context.Context) error {
	err := b.createServiceAccount(ctx, "users:user1", "user1-sa")
	if err != nil {
		return err
	}

	err = b.createServiceAccountKubeconfig(ctx, "users:user1", "user1-sa", "./dev/user1.kubeconfig")
	if err != nil {
		return err
	}

	err = b.createServiceAccountRoleBinding(ctx, "users:user1", "user1-sa")
	if err != nil {
		return err
	}

	return nil
}

func newKCPClusterClient(config *rest.Config) (kcpclient.ClusterInterface, error) {
	clusterConfig := rest.CopyConfig(config)
	u, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}
	u.Path = ""
	clusterConfig.Host = u.String()
	clusterConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	return kcpclient.NewClusterForConfig(clusterConfig)
}
