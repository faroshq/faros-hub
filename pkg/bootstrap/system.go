package bootstrap

import (
	"context"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"

	"github.com/faroshq/faros-hub/pkg/bootstrap/templates/systemtenants"
	bootstraputils "github.com/faroshq/faros-hub/pkg/util/bootstrap"
)

func (b *bootstrap) bootstrapSystemTenantAssets(ctx context.Context, workspace string) error {
	rest, err := b.clientFactory.GetWorkspaceRestConfig(ctx, workspace)
	if err != nil {
		return err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(rest)
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(rest)
	if err != nil {
		return err
	}

	return systemtenants.Bootstrap(ctx, discoveryClient, dynamicClient, bootstraputils.ReplaceOption())
}