package bootstrap

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"

	"github.com/davecgh/go-spew/spew"
	"github.com/kcp-dev/logicalcluster/v2"

	"github.com/faroshq/faros-hub/pkg/bootstrap/templates/root"
	"github.com/faroshq/faros-hub/pkg/bootstrap/templates/servicetenants"
	bootstraputils "github.com/faroshq/faros-hub/pkg/util/bootstrap"
)

func (b *bootstrap) bootstrapServiceTenantAssets(ctx context.Context, workspace string) error {
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

	ctxC := logicalcluster.WithCluster(ctx, logicalcluster.New(b.config.ControllersWorkspace))
	export, err := b.kcpClient.ApisV1alpha1().APIExports().Get(ctxC, "tenancy.faros.sh", metav1.GetOptions{})
	if err != nil {
		spew.Dump(err)
		return err
	}

	return servicetenants.Bootstrap(ctx, discoveryClient, dynamicClient, bootstraputils.ReplaceOption(
		"IDENTITY", export.Status.IdentityHash,
	))
}

func (b *bootstrap) bootstrapRootTenantAssets(ctx context.Context) error {
	rest, err := b.clientFactory.GetWorkspaceRestConfig(ctx, "root")
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

	return root.Bootstrap(ctx, discoveryClient, dynamicClient, bootstraputils.ReplaceOption())
}
