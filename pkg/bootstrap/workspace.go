package bootstrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	tenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	pluginhelpers "github.com/kcp-dev/kcp/pkg/cliplugins/helpers"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (b *bootstrap) createNamedWorkspace(ctx context.Context, workspace string) error {
	client, rest, err := b.clientFactory.GetChildWorkspaceKCPClient(ctx, workspace)
	if err != nil {
		return err
	}

	_, currentClusterName, err := pluginhelpers.ParseClusterURL(rest.Host)
	if err != nil {
		return fmt.Errorf("current URL %q does not point to cluster workspace", b.config.RootRestConfig.Host)
	}

	separatorIndex := strings.LastIndex(workspace, ":")
	var structuredWorkspaceType tenancyv1alpha1.ClusterWorkspaceTypeReference
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
