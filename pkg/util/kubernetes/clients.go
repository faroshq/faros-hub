package kubernetes

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	pluginhelpers "github.com/kcp-dev/kcp/pkg/cliplugins/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type ClientFactory interface {
	GetWorkspaceClient(ctx context.Context, workspace string) (kcpclient.ClusterInterface, *rest.Config, error)
	GetChildWorkspaceClient(ctx context.Context, workspace string) (kcpclient.ClusterInterface, *rest.Config, error)
	GetRootKCPClusterClient(ctx context.Context) (kcpclient.ClusterInterface, error)
}

type clientFactory struct {
	rest *rest.Config
}

func NewClientFactory(config *rest.Config) (*clientFactory, error) {
	return &clientFactory{
		rest: config,
	}, nil
}

func (c *clientFactory) GetRootKCPClusterClient(_ context.Context) (kcpclient.ClusterInterface, error) {
	clusterConfig := rest.CopyConfig(c.rest)
	u, err := url.Parse(c.rest.Host)
	if err != nil {
		return nil, err
	}
	u.Path = ""
	clusterConfig.Host = u.String()
	clusterConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	return kcpclient.NewClusterForConfig(clusterConfig)
}

func (c *clientFactory) GetChildWorkspaceClient(ctx context.Context, workspace string) (kcpclient.ClusterInterface, *rest.Config, error) {
	client, err := c.GetRootKCPClusterClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return c.getChildWorkspaceClient(ctx, client, c.rest, workspace)
}

func (c *clientFactory) GetWorkspaceClient(ctx context.Context, workspace string) (kcpclient.ClusterInterface, *rest.Config, error) {
	client, err := c.GetRootKCPClusterClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return c.getWorkspaceClient(ctx, client, c.rest, workspace)
}

func (c *clientFactory) getChildWorkspaceClient(ctx context.Context, client kcpclient.ClusterInterface, config *rest.Config, workspace string) (kcpclient.ClusterInterface, *rest.Config, error) {
	_, currentClusterName, err := pluginhelpers.ParseClusterURL(config.Host)
	if err != nil {
		return nil, nil, fmt.Errorf("current URL %q does not point to cluster workspace", config.Host)
	}

	parts := strings.Split(workspace, ":")

	if len(parts) >= 2 {
		currentWorkspace := parts[0]
		childWorkspace := strings.Join(parts[1:], ":")

		ws, err := client.Cluster(currentClusterName).TenancyV1beta1().Workspaces().Get(ctx, currentWorkspace, metav1.GetOptions{})
		if err != nil {
			return nil, nil, err
		}

		u, err := url.Parse(ws.Status.URL)
		if err != nil {
			return nil, nil, err
		}

		clusterConfig := rest.CopyConfig(config)
		clusterConfig.Host = u.String()
		clusterConfig.UserAgent = rest.DefaultKubernetesUserAgent()

		return c.getChildWorkspaceClient(ctx, client, clusterConfig, childWorkspace)
	}
	return client, config, nil
}

func (c *clientFactory) getWorkspaceClient(ctx context.Context, client kcpclient.ClusterInterface, config *rest.Config, workspace string) (kcpclient.ClusterInterface, *rest.Config, error) {
	client, config, err := c.getChildWorkspaceClient(ctx, client, config, workspace)
	if err != nil {
		return nil, nil, err
	}

	parts := strings.Split(workspace, ":")

	_, currentClusterName, err := pluginhelpers.ParseClusterURL(config.Host)
	if err != nil {
		return nil, nil, fmt.Errorf("current URL %q does not point to cluster workspace", config.Host)
	}

	ws, err := client.Cluster(currentClusterName).TenancyV1beta1().Workspaces().Get(ctx, parts[len(parts)-1], metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	u, err := url.Parse(ws.Status.URL)
	if err != nil {
		return nil, nil, err
	}

	clusterConfig := rest.CopyConfig(config)
	clusterConfig.Host = u.String()
	clusterConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	return client, clusterConfig, nil
}
