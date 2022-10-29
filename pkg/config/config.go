package config

import "k8s.io/client-go/rest"

const (
	ConfigFileName = "config.yaml"
)

type ControllerConfig struct {
	// ControllerExternalURL is the URL that the controller is externally reachable at.
	ControllerExternalURL string `envconfig:"FAROS_CONTROLLER_EXTERNAL_URL" required:"true" default:"https://localhost:6443"`
	// ControllersFarosEdgeAPIExportName is name of edge api export
	ControllersFarosEdgeAPIExportName string `envconfig:"FAROS_CONTROLLER_EDGE_APIEXPORT" yaml:"controllersEdgeAPIExport,omitempty" default:"edge.faros.sh"`
	// ControllersFarosPluginsAPIExportName is name of plugins api export
	ControllersFarosPluginsAPIExportName string `envconfig:"FAROS_CONTROLLER_PLUGINS_APIEXPORT" yaml:"controllersPluginsAPIExport,omitempty" default:"plugins.faros.sh"`
	// ControllersFarosAccessAPIExportName is name of access api export
	ControllersFarosAccessAPIExportName string `envconfig:"FAROS_CONTROLLER_ACCESS_APIEXPORT" yaml:"controllersAccessAPIExport,omitempty" default:"access.faros.sh"`
	// ControllersFarosAccessAPIExportName is name of access api export
	ControllersFarosTenancyAPIExportName string `envconfig:"FAROS_CONTROLLER_TENANCY_APIEXPORT" yaml:"controllersTenancyAPIExport,omitempty" default:"tenancy.faros.sh"`

	// ControllersWorkspace is name of workspace controllers are operating in
	ControllersWorkspace string `envconfig:"FAROS_CONTROLLER_WORKSPACE" yaml:"controllersWorkspace,omitempty" default:"root:faros:service:controllers"`
	// ControllersTenantWorkspace is name of workspace for global tenant management. Used in service management
	ControllersTenantWorkspace string `envconfig:"FAROS_CONTROLLER_TENANT_WORKSPACE" yaml:"controllersTenantWorkspace,omitempty" default:"root:faros:service:tenants"`

	// TenantsWorkspacePrefix is prefix of workspace tenants are operating in
	// TODO: Move under users once we can rebase to main. There is some bug in
	// using homedir but I was not able to reproduce it in main kcp branch so I am not sure if it is fixed
	TenantsWorkspacePrefix string `envconfig:"FAROS_TENANTS_WORKSPACE_PREFIX" yaml:"tenantsWorkspacePrefix,omitempty" default:"root:faros-tenants"`

	// TenantsCertificateAuthorityData is the base64 encoded certificate authority data for the tenants KubeConfigs. If not set it will set
	// skip TLS verification for the tenants KubeConfigs
	TenantsCertificateAuthorityData string `envconfig:"FAROS_TENANTS_CA_DATA" yaml:"tenantsCertificateAuthorityData,omitempty" default:""`

	// OIDC provider configuration. We will route user request to this provider and wait for callback to our API with credentials
	OIDCIssuerURL      string `envconfig:"FAROS_OIDC_ISSUER_URL" yaml:"oidcIssuerURL,omitempty" default:"https://dex.dev.faros.sh"`
	OIDCClientID       string `envconfig:"FAROS_OIDC_CLIENT_ID" yaml:"oidcClientID,omitempty" default:"faros"`
	OIDCClientSecret   string `envconfig:"FAROS_OIDC_CLIENT_SECRET" yaml:"oidcClientSecret,omitempty" default:"faros"`
	OIDCCAFile         string `envconfig:"FAROS_OIDC_CA_FILE" yaml:"oidcCAFile,omitempty" default:"hack/dev/dex/ssl/ca.pem"`
	OIDCUsernameClaim  string `envconfig:"FAROS_OIDC_USERNAME_CLAIM" yaml:"oidcFarosUsernameClaim,omitempty" default:"email"`
	OIDCUserPrefix     string `envconfig:"FAROS_OIDC_USER_PREFIX" yaml:"oidcUserPrefix,omitempty" default:"faros-sso"`
	OIDCGroupsPrefix   string `envconfig:"FAROS_OIDC_GROUPS_PREFIX" yaml:"oidcGroupsPrefix,omitempty" default:"faros-sso"`
	OIDCAuthSessionKey string `envconfig:"FAROS_OIDC_AUTH_SESSION_KEY" yaml:"oidcAuthSessionKey,omitempty" default:""`

	RestConfig *rest.Config `yaml:"-"`
}

type AgentConfig struct {
	Name      string `envconfig:"FAROS_AGENT_NAME" yaml:"name,omitempty" default:""`
	Namespace string `envconfig:"FAROS_AGENT_NAMESPACE" yaml:"namespace,omitempty" default:""`

	RestConfig *rest.Config `yaml:"-"`
}
