package kubernetes

import (
	"context"
	"net/url"

	kcpclienthelper "github.com/kcp-dev/apimachinery/pkg/client"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v3"
	"k8s.io/client-go/rest"
)

type ClientFactory interface {
	GetRootKCPClient() (kcpclient.Interface, error)

	GetRootRestConfig() (*rest.Config, error)
	GetWorkspaceRestConfig(ctx context.Context, workspace string) (*rest.Config, error)
}

type clientFactory struct {
	rest             *rest.Config
	kcpClusterClient kcpclient.Interface
}

func NewClientFactory(config *rest.Config) (*clientFactory, error) {
	client, err := kcpclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &clientFactory{
		rest:             config,
		kcpClusterClient: client,
	}, nil
}

func (c *clientFactory) GetRootKCPClient() (kcpclient.Interface, error) {
	clusterConfig, err := c.getRootRestConfig()
	if err != nil {
		return nil, err
	}
	return kcpclient.NewForConfig(clusterConfig)
}

func (c *clientFactory) GetRootRestConfig() (*rest.Config, error) {
	return c.getRootRestConfig()
}

func (c *clientFactory) GetWorkspaceRestConfig(ctx context.Context, workspace string) (*rest.Config, error) {
	clusterPath := logicalcluster.NewPath(workspace)

	u, err := url.Parse(c.rest.Host)
	if err != nil {
		return nil, err
	}

	r := rest.CopyConfig(c.rest)
	u.Path = clusterPath.RequestPath()
	r.Host = u.String()

	return kcpclienthelper.SetMultiClusterRoundTripper(rest.CopyConfig(r)), nil
}

func (c *clientFactory) getRootRestConfig() (*rest.Config, error) {
	clusterConfig := rest.CopyConfig(c.rest)
	u, err := url.Parse(c.rest.Host)
	if err != nil {
		return nil, err
	}
	u.Path = ""
	clusterConfig.Host = u.String()
	clusterConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	return kcpclienthelper.SetMultiClusterRoundTripper(rest.CopyConfig(clusterConfig)), nil
}
