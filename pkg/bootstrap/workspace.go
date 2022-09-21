package bootstrap

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	tenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	pluginhelpers "github.com/kcp-dev/kcp/pkg/cliplugins/helpers"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
)

func (b *bootstrap) createNamedWorkspace(ctx context.Context, workspace string) error {
	client, rest, err := b.getWorkspaceClient(ctx, b.kcpClient, b.rest, workspace)
	if err != nil {
		return err
	}

	_, currentClusterName, err := pluginhelpers.ParseClusterURL(rest.Host)
	if err != nil {
		return fmt.Errorf("current URL %q does not point to cluster workspace", b.rest.Host)
	}

	separatorIndex := strings.LastIndex(workspace, ":")
	var structuredWorkspaceType tenancyv1alpha1.ClusterWorkspaceTypeReference
	//switch separatorIndex {
	//case -1:
	//	structuredWorkspaceType = tenancyv1alpha1.ClusterWorkspaceTypeReference{
	//		Name: tenancyv1alpha1.ClusterWorkspaceTypeName(strings.ToLower("organization")),
	//		// path is defaulted through admission
	//	}
	//default:
	//	structuredWorkspaceType = tenancyv1alpha1.ClusterWorkspaceTypeReference{
	//		Name: tenancyv1alpha1.ClusterWorkspaceTypeName(strings.ToLower("universal")),
	//		Path: workspace[:separatorIndex],
	//	}
	//}
	spew.Dump("create in", currentClusterName, workspace, workspace[separatorIndex+1:])
	ws, err := client.Cluster(currentClusterName).TenancyV1beta1().Workspaces().Create(ctx, &tenancyv1beta1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name: workspace[separatorIndex+1:],
		},
		Spec: tenancyv1beta1.WorkspaceSpec{
			Type: structuredWorkspaceType,
		},
	}, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		ws, err = client.Cluster(currentClusterName).TenancyV1beta1().Workspaces().Get(ctx, workspace[separatorIndex+1:], metav1.GetOptions{})
	}
	if err != nil {
		return err
	}

	if err := wait.PollImmediate(time.Millisecond*100, time.Second*5, func() (bool, error) {
		if _, err := client.Cluster(currentClusterName).TenancyV1beta1().Workspaces().Get(ctx, ws.Name, metav1.GetOptions{}); err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}); err != nil {
		return err
	}

	// wait for being ready
	if ws.Status.Phase != tenancyv1alpha1.ClusterWorkspacePhaseReady {
		if err := wait.PollImmediate(time.Millisecond*500, time.Second*5, func() (bool, error) {
			ws, err = client.Cluster(currentClusterName).TenancyV1beta1().Workspaces().Get(ctx, ws.Name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if ws.Status.Phase == tenancyv1alpha1.ClusterWorkspacePhaseReady {
				return true, nil
			}
			return false, nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (b *bootstrap) getWorkspaceClient(ctx context.Context, client kcpclient.ClusterInterface, config *rest.Config, workspace string) (kcpclient.ClusterInterface, *rest.Config, error) {
	_, currentClusterName, err := pluginhelpers.ParseClusterURL(config.Host)
	if err != nil {
		return nil, nil, fmt.Errorf("current URL %q does not point to cluster workspace", b.rest.Host)
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

		return b.getWorkspaceClient(ctx, client, clusterConfig, childWorkspace)
	} else {
		u, err := url.Parse(config.Host)
		if err != nil {
			return nil, nil, err
		}
		u.Path = ""

		clusterConfig := rest.CopyConfig(config)
		clusterConfig.Host = u.String()
		clusterConfig.UserAgent = rest.DefaultKubernetesUserAgent()
		client, err = kcpclient.NewClusterForConfig(clusterConfig)
		if err != nil {
			return nil, nil, err
		}
		return client, config, nil
	}

}
