package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Load loads the configuration from the environment and flags
// Loading order:
// 1. Load .env file
// 2. Load envconfig from ENV variables and defaults
func Load() (*Config, error) {
	c := &Config{}
	// 1. Load .env file
	godotenv.Load()

	// 2. Load ENV and defaults
	err := envconfig.Process("", c)
	if err != nil {
		return c, err
	}

	// load root rest config
	restConfig := ctrl.GetConfigOrDie()
	c.RootRestConfig = restConfig

	return c, err
}
