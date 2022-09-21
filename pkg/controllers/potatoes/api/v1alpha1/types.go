package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PotatoSpec defines the desired state of Potato object
type PotatoSpec struct {
	Request int64 `json:"request,omitempty"`
}

// PotatoStatus defines the observed state of Potato object
type PotatoStatus struct {
	Total   int64  `json:"total,omitempty"`
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Potato is the Schema for the Potato API
type Potato struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PotatoSpec   `json:"spec,omitempty"`
	Status PotatoStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PotatoList contains a list of Potato
type PotatoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Potato `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Potato{}, &PotatoList{})
}
