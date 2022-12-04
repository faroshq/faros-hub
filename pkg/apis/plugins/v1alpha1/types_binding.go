package v1alpha1

import (
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +crd
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:object:root=true

// Binding is the Schema for the binding API
type Binding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BindingSpec   `json:"spec,omitempty"`
	Status BindingStatus `json:"status,omitempty"`
}

// BindingSpec defines the desired state of plugin binding
type BindingSpec struct {
	// PluginType is the type of plugin to bind to.
	PluginType metav1.TypeMeta `json:"pluginType,omitempty"`
	// PluginName is the name of the plugin to bind to.
	PluginName string `json:"pluginName,omitempty"`
	// Selector is a label selector for agents to bind to.
	Selector metav1.LabelSelector `json:"selector,omitempty"`
}

// BindingStatus defines the observed state of plugin binding
type BindingStatus struct {
	// Current processing state of the Agent.
	// +optional
	Conditions conditionsv1alpha1.Conditions `json:"conditions,omitempty"`
}

func (in *Binding) SetConditions(c conditionsv1alpha1.Conditions) {
	in.Status.Conditions = c
}

func (in *Binding) GetConditions() conditionsv1alpha1.Conditions {
	return in.Status.Conditions
}

// BindingList contains a list of bindings
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type BindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Binding `json:"items"`
}
