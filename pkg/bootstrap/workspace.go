package bootstrap

import (
	"context"
	"strings"
	"time"

	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	tenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	"github.com/kcp-dev/logicalcluster/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (b *bootstrap) createNamedWorkspace(ctx context.Context, workspace string) error {
	// create nested workspaces if requested one does not have pattern created
	separatorIndex := strings.LastIndex(workspace, ":")
	name := workspace[separatorIndex+1:]
	parent, exists := logicalcluster.New(workspace).Parent()
	ctx = logicalcluster.WithCluster(ctx, parent)

	if exists {
		_, err := b.kcpClient.TenancyV1alpha1().ClusterWorkspaces().Get(ctx, name, metav1.GetOptions{})
		if err != nil && (apierrors.IsNotFound(err) || apierrors.IsForbidden(err)) {
			switch {
			case apierrors.IsNotFound(err) && parent.String() == string(tenancyv1alpha1.RootWorkspaceTypeName):
				// ok, flow below will create root workspace
			case apierrors.IsForbidden(err) && parent.String() == string(tenancyv1alpha1.RootWorkspaceTypeName):
				return err
			case (apierrors.IsForbidden(err) && parent.String() != string(tenancyv1alpha1.RootWorkspaceTypeName)) ||
				(apierrors.IsNotFound(err) && parent.String() != string(tenancyv1alpha1.RootWorkspaceTypeName)):
				// create parent workspace if it does not exist
				err = b.createNamedWorkspace(ctx, parent.String())
				if err != nil {
					return err
				}
			default:
				return err
			}
		}
	}

	var structuredWorkspaceType tenancyv1alpha1.ClusterWorkspaceTypeReference
	ws, err := b.kcpClient.TenancyV1beta1().Workspaces().Create(ctx, &tenancyv1beta1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: tenancyv1beta1.WorkspaceSpec{
			Type: structuredWorkspaceType,
		},
	}, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	if err := wait.PollImmediate(time.Millisecond*100, time.Second*5, func() (bool, error) {
		if _, err := b.kcpClient.TenancyV1beta1().Workspaces().Get(ctx, name, metav1.GetOptions{}); err != nil {
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
			ws, err = b.kcpClient.TenancyV1beta1().Workspaces().Get(ctx, name, metav1.GetOptions{})
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
