package bootstrap

import (
	"context"
	"strings"
	"time"

	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	tenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (b *bootstrap) createNamedWorkspace(ctx context.Context, workspace string) error {
	rest, err := b.clientFactory.GetChildWorkspaceRestConfig(ctx, workspace)
	if err != nil {
		return err
	}

	client, err := kcpclient.NewForConfig(rest)
	if err != nil {
		return err
	}

	separatorIndex := strings.LastIndex(workspace, ":")
	var structuredWorkspaceType tenancyv1alpha1.ClusterWorkspaceTypeReference
	ws, err := client.TenancyV1beta1().Workspaces().Create(ctx, &tenancyv1beta1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name: workspace[separatorIndex+1:],
		},
		Spec: tenancyv1beta1.WorkspaceSpec{
			Type: structuredWorkspaceType,
		},
	}, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		ws, err = client.TenancyV1beta1().Workspaces().Get(ctx, workspace[separatorIndex+1:], metav1.GetOptions{})
	}
	if err != nil {
		return err
	}

	if err := wait.PollImmediate(time.Millisecond*100, time.Second*5, func() (bool, error) {
		if _, err := client.TenancyV1beta1().Workspaces().Get(ctx, ws.Name, metav1.GetOptions{}); err != nil {
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
			ws, err = client.TenancyV1beta1().Workspaces().Get(ctx, ws.Name, metav1.GetOptions{})
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
