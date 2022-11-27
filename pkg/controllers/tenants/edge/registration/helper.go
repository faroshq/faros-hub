package registration

import (
	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getRegistrationResourceName(registration *edgev1alpha1.Registration) string {
	return "registration-" + registration.Name
}

func getRegistrationOwnersReference(registration *edgev1alpha1.Registration) []metav1.OwnerReference {
	return []metav1.OwnerReference{{
		APIVersion: edgev1alpha1.SchemeGroupVersion.String(),
		Kind:       edgev1alpha1.RegistrationKind,
		Name:       registration.Name,
		UID:        registration.UID,
	}}
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
