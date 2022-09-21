package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Load loads the configuration from the environment and flags
// Loading order:
// 1. Load .env file
// 2. Load envconfig from ENV variables and defaults
func Load(loadCertificates bool) (*Config, error) {
	c := &Config{}
	// 1. Load .env file
	godotenv.Load()

	// 2. Load ENV and defaults
	err := envconfig.Process("", c)
	if err != nil {
		return c, err
	}

	return c, err
}
