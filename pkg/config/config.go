package config

import "k8s.io/client-go/rest"

const (
	ConfigFileName = "config.yaml"
)

type Config struct {
	ControllersAPIExport string `envconfig:"FAROS_CONTROLLER_APIEXPORT" yaml:"controllersAPIExport,omitempty" default:"faros.sh"`
	ControllersWorkspace string `envconfig:"FAROS_CONTROLLER_WORKSPACE" yaml:"controllersWorkspace,omitempty" default:"root:compute:controllers"`

	RootRestConfig *rest.Config `yaml:"-"`
}
