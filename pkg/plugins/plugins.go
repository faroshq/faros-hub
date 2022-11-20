package plugins

import (
	"context"

	"k8s.io/client-go/rest"
)

// Interface is the interface that plugins must implement.
// TODO: move to faros-hub repo
type Interface interface {
	// Name returns the name of the plugin.
	Name(ctx context.Context) (string, error)
	// Init initializes the plugin.
	Init(ctx context.Context, name, namespace string, config *rest.Config) error
	// Run runs the plugin.
	Run(ctx context.Context) error
	// Stop stops the plugin.
	Stop(ctx context.Context) error
}
