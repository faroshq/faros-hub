package config

import "k8s.io/client-go/rest"

const (
	ConfigFileName = "config.yaml"
)

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
	ControllersWorkspace string `envconfig:"FAROS_CONTROLLER_WORKSPACE" yaml:"controllersWorkspace,omitempty" default:"root:faros-system:controllers"`
	// ControllersTenantWorkspace is name of workspace for global tenant management. Used in system management
	ControllersTenantWorkspace string `envconfig:"FAROS_CONTROLLER_TENANT_WORKSPACE" yaml:"controllersTenantWorkspace,omitempty" default:"root:faros-system:tenants"`

	// TenantsWorkspacePrefix is prefix of workspace tenants are operating in
	TenantsWorkspacePrefix string `envconfig:"FAROS_TENANTS_WORKSPACE_PREFIX" yaml:"tenantsWorkspacePrefix,omitempty" default:"root:faros"`

	RestConfig *rest.Config `yaml:"-"`
}

type AgentConfig struct {
	Name      string `envconfig:"FAROS_AGENT_NAME" yaml:"name,omitempty" default:""`
	Namespace string `envconfig:"FAROS_AGENT_NAMESPACE" yaml:"namespace,omitempty" default:""`

	RestConfig *rest.Config `yaml:"-"`
}
