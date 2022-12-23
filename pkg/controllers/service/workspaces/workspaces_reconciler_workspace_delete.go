package workspaces

import (
	"context"
	"fmt"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
)

type kcpWorkspaceDeleteReconciler struct {
	deleteFarosWorkspace func(ctx context.Context, cluster logicalcluster.Path, workspace *tenancyv1alpha1.Workspace) error
	getWorkspaceName     func(w *tenancyv1alpha1.Workspace) string
}

func (r *kcpWorkspaceDeleteReconciler) reconcile(ctx context.Context, workspace *tenancyv1alpha1.Workspace) (reconcileStatus, error) {
	workspacePath := r.getWorkspaceName(workspace)
	clusterPath := logicalcluster.NewPath(workspacePath)

	parent, exits := clusterPath.Parent()
	if !exits {
		return reconcileStatusError, fmt.Errorf("parent cluster not found")
	}

	// delete faros workspaces in the child clusters
	err := r.deleteFarosWorkspace(ctx, parent, workspace)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	return reconcileStatusContinue, nil
}
