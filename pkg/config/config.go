package config

const (
	ConfigFileName = "config.yaml"
)

type Config struct {
	Server Server `yaml:"server,omitempty"`
}

type Server struct {
	StateDir string `envconfig:"SERVER_STATE_DIR" yaml:"stateDir,omitempty"`

	WorkspacesList []string `envconfig:"SERVER_WORKSPACES_LIST" yaml:"workspacesList,omitempty" default:"corp,corp:services,corp:services,corp:services:warehouse,corp:compute,corp:compute:services,corp:compute:shared,users:user1,users:user2,users:user3"`
	// ComputeServicesKubeconfigs is list of paths to kubeconfigs for services compute clusters
	ComputeServicesKubeconfigs []string `envconfig:"SERVER_COMPUTE_SERVICES_KUBECONFIGS" yaml:"computeServicesKubeconfigs,omitempty"`
	ComputeServiceWorkspace    string   `envconfig:"SERVER_COMPUTE_SERVICE_WORKSPACE" yaml:"computeServiceWorkspace,omitempty" default:"corp:compute:services"`

	// ComputeSharedKubeconfigs is list of paths to kubeconfigs for shared compute clusters
	ComputeSharedKubeconfigs []string `envconfig:"SERVER_COMPUTE_SHARED_KUBECONFIGS" yaml:"computeSharedKubeconfigs,omitempty"`
	ComputeSharedWorkspace   string   `envconfig:"SERVER_COMPUTE_SHARED_WORKSPACE" yaml:"computeSharedWorkspace,omitempty" default:"corp:compute:shared"`

	ComputeSyncerImage      string `envconfig:"SERVER_COMPUTE_SYNCER_IMAGE" yaml:"computeSyncerImage,omitempty" default:"quay.io/mangirdas/syncer-c2e3073d5026a8f7f2c47a50c16bdbec@sha256:9e7c1455769adbacb139054f1fe4eee10ebdbd55b6b8368cb83d3214fa6ddd54"`
	ComputeSyncerNamespace  string `envconfig:"SERVER_COMPUTE_SYNCER_NAMESPACE" yaml:"computeSyncerNamespace,omitempty" default:"default"`
	ComputeSycnerTargetName string `envconfig:"SERVER_COMPUTE_SYNCER_TARGET_NAME" yaml:"computeSyncerTargetName,omitempty" default:""`
	ComputeSyncerIDPPrefix  string `envconfig:"SERVER_COMPUTE_SYNCER_IDP_PREFIX" yaml:"computeSyncerIDPPrefix,omitempty" default:"kcp-syncer-"`

	ComputeSyncerResourcesToSync []string `envconfig:"SERVER_COMPUTE_SYNCER_RESOURCES_TO_SYNC" yaml:"computeSyncerResourcesToSync,omitempty" default:"services,ingresses.networking.k8s.io"`
	ComputeSyncerFeatureGates    string   `envconfig:"SERVER_COMPUTE_SYNCER_FEATURE_GATES" yaml:"computeSyncerFeatureGates,omitempty" default:""`

	ControllerPotatoesCount     int64  `envconfig:"SERVER_CONTROLLER_POTATOES_COUNT" yaml:"controllerPotatoesCount,omitempty" default:"100"`
	ControllerServicesWorkspace string `envconfig:"SERVER_CONTROLLER_SERVICES_WORKSPACE" yaml:"controllerServicesWorkspace,omitempty" default:"corp:services:warehouse"`
}
