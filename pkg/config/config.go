package config

import "k8s.io/client-go/rest"

const (
	ConfigFileName = "config.yaml"
)

type APIConfig struct {
	// Addr is the address to bind the controller to.
	Addr string `envconfig:"FAROS_API_ADDR" required:"true" default:":8080"`
	// ControllerExternalURL is the URL that the controller is externally reachable at.
	ControllerExternalURL string `envconfig:"FAROS_API_EXTERNAL_URL" required:"true" default:"https://kcp.dev.faros.sh"`

	// Important: HostingClusterKubeConfigPath is used to dynamically read secrets for trust. For now single secrets we
	// require in API server context is OIDC CA bundle from Dex. If removed this dependency, this can be
	// removed.
	// HostingClusterKubeConfig is the path to the kubeconfig file for the hosting cluster.
	HostingClusterKubeConfigPath string `envconfig:"FAROS_API_HOSTING_CLUSTER_KUBECONFIG" required:"true" default:"faros.kubeconfig"`
	// HostingClusterNamespace is the namespace in the hosting cluster where the controller will run.
	HostingClusterNamespace string `envconfig:"FAROS_API_HOSTING_CLUSTER_NAMESPACE" required:"true" default:"kcp"`
	// HostingClusterRestConfig is the rest config for the hosting cluster.
	// Lodaded from HostingClusterKubeConfig.
	HostingClusterRestConfig *rest.Config `envconfig:"-"`

	// KCPClusterKubeConfigPath is the path to the kubeconfig file for the kcp cluster
	KCPClusterKubeConfigPath string `envconfig:"FAROS_API_KCP_CLUSTER_KUBECONFIG" required:"true" default:"kcp.kubeconfig"`
	// KCPClusterRestConfig is the rest config for the KCP cluster.
	// Used to manage users, workspaces, etc
	KCPClusterRestConfig *rest.Config `envconfig:"-"`

	// ControllersTenantWorkspace is name of workspace for global tenant management. Used in service management
	// Must match one in Controllers config
	ControllersTenantWorkspace string `envconfig:"FAROS_API_TENANT_WORKSPACE" yaml:"controllersTenantWorkspace,omitempty" default:"root:faros:service:tenants"`

	// OIDC provider configuration
	OIDCIssuerURL      string `envconfig:"FAROS_OIDC_ISSUER_URL" yaml:"oidcIssuerURL,omitempty" default:"https://dex.dev.faros.sh"`
	OIDCClientID       string `envconfig:"FAROS_OIDC_CLIENT_ID" yaml:"oidcClientID,omitempty" default:"faros"`
	OIDCClientSecret   string `envconfig:"FAROS_OIDC_CLIENT_SECRET" yaml:"oidcClientSecret,omitempty" default:"faros"`
	OIDCCASecretName   string `envconfig:"FAROS_OIDC_CA_SECRET_NAME" yaml:"oidcCASecretName,omitempty" default:"dex-pki-ca"`
	OIDCUsernameClaim  string `envconfig:"FAROS_OIDC_USERNAME_CLAIM" yaml:"oidcFarosUsernameClaim,omitempty" default:"email"`
	OIDCUserPrefix     string `envconfig:"FAROS_OIDC_USER_PREFIX" yaml:"oidcUserPrefix,omitempty" default:"faros-sso"`
	OIDCGroupsPrefix   string `envconfig:"FAROS_OIDC_GROUPS_PREFIX" yaml:"oidcGroupsPrefix,omitempty" default:"faros-sso"`
	OIDCAuthSessionKey string `envconfig:"FAROS_OIDC_AUTH_SESSION_KEY" yaml:"oidcAuthSessionKey,omitempty" default:""`
}

type ControllerConfig struct {
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

	// TenantsCertificateAuthorityFile is the file for certificate for the tenants KubeConfigs. If not set it will set
	// skip TLS verification for the tenants KubeConfigs
	TenantsCertificateAuthorityFile string `envconfig:"FAROS_TENANTS_CA_FILE" yaml:"tenantsCertificateAuthorityFile,omitempty" default:""`

	// TenantsCertificateAuthorityData is the data for certificate for the tenants KubeConfigs. It will be set from TenantsCertificateAuthorityFile
	TenantsCertificateAuthorityFileData []byte `yaml:"tenantsCertificateAuthorityFileData,omitempty"`

	// Must match one in API config
	OIDCUsernameClaim string `envconfig:"FAROS_OIDC_USERNAME_CLAIM" yaml:"oidcFarosUsernameClaim,omitempty" default:"email"`
	OIDCUserPrefix    string `envconfig:"FAROS_OIDC_USER_PREFIX" yaml:"oidcUserPrefix,omitempty" default:"faros-sso"`

	// KCPClusterKubeConfigPath is the path to the kubeconfig file for the kcp cluster
	KCPClusterKubeConfigPath string `envconfig:"FAROS_CONTROLLER_KCP_CLUSTER_KUBECONFIG" required:"true" default:"kcp.kubeconfig"`
	// KCPClusterRestConfig is the rest config for the KCP cluster.
	KCPClusterRestConfig *rest.Config `envconfig:"-"`
}

type AgentConfig struct {
	Name      string `envconfig:"FAROS_AGENT_NAME" yaml:"name,omitempty" default:""`
	Namespace string `envconfig:"FAROS_AGENT_NAMESPACE" yaml:"namespace,omitempty" default:""`

	RestConfig *rest.Config `yaml:"-"`
}
