package config

const (
	ConfigFileName = "config.yaml"
)

type Config struct {
	Server Server `yaml:"server,omitempty"`
}

type Server struct {
	StateDir string `envconfig:"SERVER_STATE_DIR" yaml:"stateDir,omitempty"`
}
