package bootstrap

import (
	"context"

	"github.com/faroshq/faros-hub/pkg/config"
	"k8s.io/client-go/rest"

	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
)

type Bootstraper interface {
	CreateWorkspace(ctx context.Context, name string) error
	DeployKustomizeAssets(ctx context.Context, workspace string, path string) error
}

type bootstrap struct {
	rest   *rest.Config
	config *config.Config

	clientFactory utilkubernetes.ClientFactory
}

func New(config *config.Config, rest *rest.Config) (*bootstrap, error) {
	cf, err := utilkubernetes.NewClientFactory(rest)
	if err != nil {
		return nil, err
	}

	b := &bootstrap{
		rest:          rest,
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
