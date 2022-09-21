package bootstrap

import (
	"context"
	"net/url"

	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
)

type Bootstraper interface {
	BootstrapOrganization(ctx context.Context) error
	BootstrapCompute(ctx context.Context) error
}

type bootstrap struct {
	rest *rest.Config

	kcpClient kcpclient.ClusterInterface
}

func New(ctx context.Context, rest *rest.Config) (*bootstrap, error) {
	kcpClient, err := newKCPClusterClient(rest)
	if err != nil {
		return nil, err
	}

	b := &bootstrap{
		rest:      rest,
		kcpClient: kcpClient,
	}

	return b, nil
}

func (b *bootstrap) BootstrapOrganization(ctx context.Context) error {
	workspaces := []string{
		"corp",
		"corp:services",
		"corp:services:potatoes",
		"corp:services:bananas",
		"corp:compute",
		"corp:compute:services",
		"corp:compute:shared",
	}
	for _, workspace := range workspaces {
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
