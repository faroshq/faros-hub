package workspaces

import (
	"context"
	"fmt"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
)

type kcpWorkspaceDeleteReconciler struct {
	deleteFarosWorkspace func(ctx context.Context, cluster logicalcluster.Name, workspace *tenancyv1alpha1.Workspace) error
	getWorkspaceName     func(w *tenancyv1alpha1.Workspace) string
}

func (r *kcpWorkspaceDeleteReconciler) reconcile(ctx context.Context, workspace *tenancyv1alpha1.Workspace) (reconcileStatus, error) {
	workspacePath := r.getWorkspaceName(workspace)
	cluster := logicalcluster.New(workspacePath)

	parent, exits := cluster.Parent()
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
