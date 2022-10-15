package bootstrap

import (
	"context"

	"github.com/faroshq/faros-hub/pkg/config"
	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
)

type Bootstraper interface {
	CreateWorkspace(ctx context.Context, name string) error
	DeployKustomizeAssets(ctx context.Context, workspace string, path string) error
}

type bootstrap struct {
	config *config.ControllerConfig

	clientFactory utilkubernetes.ClientFactory
}

func New(config *config.ControllerConfig) (*bootstrap, error) {
	cf, err := utilkubernetes.NewClientFactory(config.RestConfig)
	if err != nil {
		return nil, err
	}

	b := &bootstrap{
		config:        config,
		clientFactory: cf,
	}

	return b, nil
}

func (b *bootstrap) DeployKustomizeAssets(ctx context.Context, workspace string, path string) error {
	err := b.deployComponents(ctx, workspace, path)
	if err != nil {
		return err
	}

	return nil
}

func (b *bootstrap) CreateWorkspace(ctx context.Context, name string) error {
	return b.createNamedWorkspace(ctx, name)
}
