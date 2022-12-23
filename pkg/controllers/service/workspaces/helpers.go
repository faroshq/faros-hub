package workspaces

import (
	"fmt"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
	kcptenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func getWorkspaceName(config *config.ControllerConfig, w *tenancyv1alpha1.Workspace) string {
	return fmt.Sprintf("%s:%s", config.TenantsWorkspacePrefix, w.Name)
}

func getOrgClusterAccessName(config *config.ControllerConfig, workspace *tenancyv1alpha1.Workspace) string {
	return fmt.Sprintf("%s-cluster-admin", workspace.Name)
}

func getUserWithPrefixName(config *config.ControllerConfig, user string) string {
	return fmt.Sprintf("%s:%s", config.OIDCUserPrefix, user)
}

func getRBACClusterAdminName(workspace *tenancyv1alpha1.Workspace) string {
	return fmt.Sprintf("%s-cluster-admin", workspace.Name)
}

// mergeOwnerReference: merge a slice of ownerReference with a given ownerReferences
func mergeOwnerReference(ownerReferences, newOwnerReferences []metav1.OwnerReference) []metav1.OwnerReference {
	var merged []metav1.OwnerReference

	merged = append(merged, ownerReferences...)

	for _, ownerReference := range newOwnerReferences {
		found := false
		for _, mergedOwnerReference := range merged {
			if mergedOwnerReference.UID == ownerReference.UID {
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, ownerReference)
		}
	}

	return merged
}

// getKCPWorkspaceOwnersReference returns the owners reference of given KCP workspace
func getKCPWorkspaceOwnersReference(workspace *kcptenancyv1beta1.Workspace) []metav1.OwnerReference {
	return []metav1.OwnerReference{{
		APIVersion:         kcptenancyv1beta1.SchemeGroupVersion.String(),
		Kind:               "Workspace",
		Name:               workspace.Name,
		BlockOwnerDeletion: pointer.BoolPtr(true),
		Controller:         pointer.BoolPtr(true),
		UID:                workspace.UID,
	}}
}
