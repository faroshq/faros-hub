package config

const (
	ConfigFileName = "config.yaml"
)

type Config struct {
	ControllersAPIExport string `envconfig:"FAROS_CONTROLLER_APIEXPORT" yaml:"controllersAPIExport,omitempty" default:"access.faros.sh"`
	ControllersWorkspace string `envconfig:"FAROS_CONTROLLER_WORKSPACE" yaml:"controllersWorkspace,omitempty" default:"compute:controllers"`
}
