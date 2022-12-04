package v1alpha1

import (
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Plugin is the generic Schema for the plugins API
// Individual plugins should have their own types. This type is used for
// listing and showing generic plugins view. It should be extended to include
// all necessary fields.

// +crd
// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:object:root=true

// Plugin is the Schema for the plugin API
type Plugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PluginSpec   `json:"spec,omitempty"`
	Status PluginStatus `json:"status,omitempty"`
}

// PluginSpec defines the desired state of plugin
type PluginSpec struct {
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

// PluginStatus defines the observed state of plugin
type PluginStatus struct {
	// Current processing state of the Agent.
	// +optional
	Conditions conditionsv1alpha1.Conditions `json:"conditions,omitempty"`
}

func (in *Plugin) SetConditions(c conditionsv1alpha1.Conditions) {
	in.Status.Conditions = c
}

func (in *Plugin) GetConditions() conditionsv1alpha1.Conditions {
	return in.Status.Conditions
}

// PluginList contains a list of plugins
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type PluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Plugin `json:"items"`
}
