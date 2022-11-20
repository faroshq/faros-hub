package bootstrap

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	apisv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	utilplugins "github.com/faroshq/faros-hub/pkg/util/plugins"
)

func (b *bootstrap) LoadPlugins(ctx context.Context, workspace string) error {
	path := b.config.PluginsDir

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	cluster := logicalcluster.New(b.config.ControllersPluginsWorkspace)

	for _, file := range files {
		p, err := utilplugins.Load(filepath.Join(path, file.Name()))
		if err != nil {
			return err
		}

		name, err := p.GetName(ctx)
		if err != nil {
			return fmt.Errorf("failed to get plugin name: %w", err)
		}

		data, err := p.GetAPIResourceSchema(ctx)
		if err != nil {
			return fmt.Errorf("failed to get API resource schema for plugin %q: %w", name, err)
		}

		var schema apisv1alpha1.APIResourceSchema
		err = yaml.Unmarshal(data, &schema)
		if err != nil {
			return fmt.Errorf("failed to unmarshal API resource schema for plugin %q: %w", name, err)
		}

		_, err = b.kcpClient.Cluster(cluster).ApisV1alpha1().APIResourceSchemas().Create(ctx, &schema, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create API resource schema for plugin %q: %w", name, err)
		}
	}

	return nil
}
