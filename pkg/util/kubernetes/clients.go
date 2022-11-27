package kubernetes

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	kcpclienthelper "github.com/kcp-dev/apimachinery/pkg/client"
	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/client-go/rest"
)

type ClientFactory interface {
	GetRootKCPClient() (kcpclient.Interface, error)

	GetRootRestConfig() (*rest.Config, error)
	GetWorkspaceRestConfig(ctx context.Context, workspace string) (*rest.Config, error)
	GetChildWorkspaceRestConfig(ctx context.Context, workspace string) (*rest.Config, error)
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
	cluster := logicalcluster.New(workspace)

	if strings.Contains(workspace, ":") && !cluster.HasPrefix(logicalcluster.New("system")) &&
		!cluster.HasPrefix(tenancyv1alpha1.RootCluster) {
		return nil, fmt.Errorf("invalid workspace name format: %s", workspace)
	}

	u, err := url.Parse(c.rest.Host)
	if err != nil {
		return nil, err
	}

	r := rest.CopyConfig(c.rest)
	u.Path = cluster.Path()
	r.Host = u.String()

	return kcpclienthelper.SetMultiClusterRoundTripper(rest.CopyConfig(r)), nil
}

func (c *clientFactory) GetChildWorkspaceRestConfig(ctx context.Context, workspace string) (*rest.Config, error) {
	parent, exists := logicalcluster.New(workspace).Parent()
	if !exists {
		return nil, fmt.Errorf("workspace %q  does not have child workspace", workspace)
	}

	if strings.Contains(workspace, ":") && !parent.HasPrefix(logicalcluster.New("system")) &&
		!parent.HasPrefix(tenancyv1alpha1.RootCluster) {
		return nil, fmt.Errorf("invalid workspace name format: %s", workspace)
	}

	u, err := url.Parse(c.rest.Host)
	if err != nil {
		return nil, err
	}

	r := rest.CopyConfig(c.rest)
	u.Path = parent.Path()
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
