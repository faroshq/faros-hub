package v1alpha1

import (
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AgentReady conditionsv1alpha1.ConditionType = "Ready"
)

// +crd
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:object:root=true

// Agent is the Schema for the Agent API
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentSpec   `json:"spec,omitempty"`
	Status AgentStatus `json:"status,omitempty"`
}

// AgentSpec defines the desired state of agent
type AgentSpec struct{}

// AgentStatus defines the observed state of Agent
type AgentStatus struct {
	// Current processing state of the Agent.
	// +optional
	Conditions conditionsv1alpha1.Conditions `json:"conditions,omitempty"`
}

func (in *Agent) SetConditions(c conditionsv1alpha1.Conditions) {
	in.Status.Conditions = c
}

func (in *Agent) GetConditions() conditionsv1alpha1.Conditions {
	return in.Status.Conditions
}

// AgentList contains a list of Agents
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}
