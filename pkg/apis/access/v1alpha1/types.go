package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KubeConfigRequest is the Schema for the KubeConfigRequest API
type KubeConfigRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeConfigRequestSpec   `json:"spec,omitempty"`
	Status KubeConfigRequestStatus `json:"status,omitempty"`
}

// KubeConfigRequestSpec defines the desired state of KubeConfigRequest
type KubeConfigRequestSpec struct {
	Request int64 `json:"request,omitempty"`
}

// PotatoStatus defines the observed state of Potato object
type KubeConfigRequestStatus struct {
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true

// KubeConfigRequestList contains a list of KubeConfigRequest
type KubeConfigRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeConfigRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeConfigRequest{}, &KubeConfigRequestList{})
}
