package workspaces

import (
	"fmt"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

func (r *Reconciler) getWorkspaceName(w *tenancyv1alpha1.Workspace) string {
	return fmt.Sprintf("%s:%s:%s", r.Config.TenantsWorkspacePrefix, w.Namespace, w.Name)
}

func (r *Reconciler) getRBACClusterAdminName(workspace *tenancyv1alpha1.Workspace) string {
	return fmt.Sprintf("%s-cluster-admin", workspace.Name)
}

func (r *Reconciler) getOrgClusterAccessName(workspace *tenancyv1alpha1.Workspace) string {
	return fmt.Sprintf("%s-%s-cluster-admin", workspace.Namespace, workspace.Name)
}
