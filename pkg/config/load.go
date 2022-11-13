package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
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

	if c.OIDCAuthSessionKey == "" {
		fmt.Println("FAROS_OIDC_AUTH_SESSION_KEY not supplied, generating random one")
		c.OIDCAuthSessionKey = uuid.Must(uuid.NewUUID()).String()
	}

	if _, err := os.Stat(c.TenantsCertificateAuthorityFile); os.IsNotExist(err) {
		return c, fmt.Errorf("tenants certificate authority file %s does not exist", c.TenantsCertificateAuthorityFile)
	}

	caServerCert, err := ioutil.ReadFile(c.TenantsCertificateAuthorityFile)
	if err != nil {
		return c, fmt.Errorf("failed to read tenants certificate authority file %s", c.TenantsCertificateAuthorityFile)
	}

	c.TenantsCertificateAuthorityFileData = caServerCert

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
