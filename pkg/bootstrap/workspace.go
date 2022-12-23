package bootstrap

import (
	"context"
	"time"

	corev1alpha1 "github.com/kcp-dev/kcp/pkg/apis/core/v1alpha1"
	tenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	"github.com/kcp-dev/logicalcluster/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (b *bootstrap) createNamedWorkspace(ctx context.Context, workspace string) error {
	clusterPath := logicalcluster.NewPath(workspace)

	_, err := b.kcpClient.Cluster(clusterPath).TenancyV1alpha1().ClusterWorkspaces().Get(ctx, workspace, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	var structuredWorkspaceType tenancyv1beta1.WorkspaceTypeReference
	ws, err := b.kcpClient.Cluster(clusterPath).TenancyV1beta1().Workspaces().Create(ctx, &tenancyv1beta1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name: workspace,
		},
		Spec: tenancyv1beta1.WorkspaceSpec{
			Type: structuredWorkspaceType,
		},
	}, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	if err := wait.PollImmediate(time.Millisecond*100, time.Second*5, func() (bool, error) {
		if _, err := b.kcpClient.Cluster(clusterPath).TenancyV1beta1().Workspaces().Get(ctx, workspace, metav1.GetOptions{}); err != nil {
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
	if ws.Status.Phase != corev1alpha1.LogicalClusterPhaseReady {
		if err := wait.PollImmediate(time.Millisecond*500, time.Second*5, func() (bool, error) {
			ws, err = b.kcpClient.Cluster(clusterPath).TenancyV1beta1().Workspaces().Get(ctx, workspace, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if ws.Status.Phase == corev1alpha1.LogicalClusterPhaseReady {
				return true, nil
			}
			return false, nil
		}); err != nil {
			return err
		}
	}

	return nil
}
