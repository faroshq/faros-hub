package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// All this is fake kube constructs to keep the consistency with the rest of the codebase

// GroupName
var GroupName = "config.faros.sh"

// CLIConfigKind is the kind for a CLIConfig
const CLIConfigKind = "CLIConfig"

type CLIConfig struct {
	metav1.TypeMeta   `json:",inline",yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec CLIConfigSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type CLIConfigSpec struct {
	BaseURL string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`
	Token   string `json:"token,omitempty" yaml:"token,omitempty"`
	Email   string `json:"email,omitempty" yaml:"email,omitempty"`
}

func NewCLIConfig() *CLIConfig {
	return &CLIConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       CLIConfigKind,
			APIVersion: GroupName + "/v1alpha1",
		},
	}
}
