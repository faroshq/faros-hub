package pluginrequests

import (
	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func getWorkspaceOwnersReference(user *tenancyv1alpha1.User) []metav1.OwnerReference {
	return []metav1.OwnerReference{{
		APIVersion:         tenancyv1alpha1.SchemeGroupVersion.String(),
		Kind:               tenancyv1alpha1.WorkspaceKind,
		Name:               user.Name,
		BlockOwnerDeletion: pointer.BoolPtr(true),
		Controller:         pointer.BoolPtr(true),
		UID:                user.UID,
	}}
}
