package workspaces

import (
	"fmt"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

func (r *Reconciler) getWorkspaceName(w *tenancyv1alpha1.Workspace) string {
	return fmt.Sprintf("%s:%s:%s", r.Config.TenantsWorkspacePrefix, w.Namespace, w.Name)
}
