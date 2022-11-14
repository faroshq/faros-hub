package bootstrap

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/faroshq/faros-hub/pkg/config"
	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
)

// TODO: All this package should go away once we have a proper bootstrap
// mechanism in place. For now, we just deploy the components and create the
// resources. Once resources are stable, we should move it to a proper code.

//go:generate go run github.com/go-bindata/go-bindata/v3/go-bindata -pkg $GOPACKAGE -prefix ../../config/ -nometadata -o zz_$GOPACKAGE.go ../../config/...
//go:generate go run golang.org/x/tools/cmd/goimports -local github.com/faroshq/faros-hub -e -w zz_$GOPACKAGE.go

type Bootstraper interface {
	CreateWorkspace(ctx context.Context, name string) error
	BootstrapServiceTenantAssets(ctx context.Context, workspace string) error
	DeployKustomizeAssetsCRD(ctx context.Context, workspace string) error
	DeployKustomizeAssetsKCP(ctx context.Context, workspace string) error
}

type bootstrap struct {
	config *config.ControllerConfig

	kcpClient     kcpclient.ClusterInterface
	clientFactory utilkubernetes.ClientFactory
}

func New(config *config.ControllerConfig) (*bootstrap, error) {
	cf, err := utilkubernetes.NewClientFactory(config.KCPClusterRestConfig)
	if err != nil {
		return nil, err
	}

	rootRest, err := cf.GetRootRestConfig()
	if err != nil {
		return nil, err
	}

	client, err := kcpclient.NewClusterForConfig(rootRest)
	if err != nil {
		return nil, err
	}

	b := &bootstrap{
		config:        config,
		clientFactory: cf,
		kcpClient:     client,
	}

	return b, nil
}

func (b *bootstrap) DeployKustomizeAssetsCRD(ctx context.Context, workspace string) error {
	tmpDir, err := os.MkdirTemp("", "faros-crd")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	for _, name := range AssetNames() {
		if strings.HasPrefix(name, "crds") {
			data, err := Asset(name)
			if err != nil {
				return err
			}
			dir := filepath.Dir(name)
			if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(tmpDir, name), data, 0644); err != nil {
				return err
			}
		}
	}

	err = b.deployComponents(ctx, workspace, tmpDir+"/crds")
	if err != nil {
		return err
	}

	return nil
}

func (b *bootstrap) DeployKustomizeAssetsKCP(ctx context.Context, workspace string) error {
	tmpDir, err := os.MkdirTemp("", "faros-kcp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	for _, name := range AssetNames() {
		if strings.HasPrefix(name, "kcp") {
			data, err := Asset(name)
			if err != nil {
				return err
			}
			dir := filepath.Dir(name)
			if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(tmpDir, name), data, 0644); err != nil {
				return err
			}
		}
	}

	err = b.deployComponents(ctx, workspace, tmpDir+"/kcp")
	if err != nil {
		return err
	}
	return nil
}

func (b *bootstrap) CreateWorkspace(ctx context.Context, name string) error {
	return b.createNamedWorkspace(ctx, name)
}

func (b *bootstrap) BootstrapServiceTenantAssets(ctx context.Context, workspace string) error {
	return b.bootstrapServiceTenantAssets(ctx, workspace)
}
