package config

import "k8s.io/client-go/rest"

const (
	ConfigFileName = "config.yaml"
)

type ControllerConfig struct {
	ControllersAPIExport string `envconfig:"FAROS_CONTROLLER_APIEXPORT" yaml:"controllersAPIExport,omitempty" default:"faros.sh"`
	ControllersWorkspace string `envconfig:"FAROS_CONTROLLER_WORKSPACE" yaml:"controllersWorkspace,omitempty" default:"root:compute:controllers"`

	RestConfig *rest.Config `yaml:"-"`
}

type AgentConfig struct {
	Name      string `envconfig:"FAROS_AGENT_NAME" yaml:"name,omitempty" default:""`
	Namespace string `envconfig:"FAROS_AGENT_NAMESPACE" yaml:"namespace,omitempty" default:""`

	RestConfig *rest.Config `yaml:"-"`
}
