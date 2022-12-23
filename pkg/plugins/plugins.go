package plugins

import (
	"context"

	"k8s.io/client-go/rest"
)

// Interface is the interface that plugins must implement.
// TODO: move to faros-hub repo
type Interface interface {
	// Name returns the name of the plugin.
	GetName(ctx context.Context) (string, error)
	// GetVersion returns the version of the plugin.
	GetVersion(context.Context) (string, error)
	// GetAPIResourceSchema returns the schema of the API resource.
	GetAPIResourceSchema(ctx context.Context) ([]byte, error)
	// GetAPIExportSchema returns the schema of the API export.
	GetAPIExportSchema(ctx context.Context) ([]byte, error)
	// Init initializes the plugin.
	Init(ctx context.Context, name, namespace string, config *rest.Config) error
	// Run runs the plugin.
	Run(ctx context.Context) error
	// Stop stops the plugin.
	Stop(ctx context.Context) error
}
