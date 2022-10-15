package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	ctrl "sigs.k8s.io/controller-runtime"
)

// LoadController loads the configuration from the environment and flags
// Loading order:
// 1. Load .env file
// 2. Load envconfig from ENV variables and defaults
func LoadController() (*ControllerConfig, error) {
	c := &ControllerConfig{}
	godotenv.Load()

	err := envconfig.Process("", c)
	if err != nil {
		return c, err
	}

	// load root rest config
	restConfig := ctrl.GetConfigOrDie()
	c.RestConfig = restConfig

	return c, err
}

// LoadAgent loads the configuration from the environment and flags
// Loading order:
// 1. Load .env file
// 2. Load envconfig from ENV variables and defaults
func LoadAgent() (*AgentConfig, error) {
	c := &AgentConfig{}
	godotenv.Load()

	err := envconfig.Process("", c)
	if err != nil {
		return c, err
	}

	// load root rest config
	restConfig := ctrl.GetConfigOrDie()
	c.RestConfig = restConfig

	return c, err
}
