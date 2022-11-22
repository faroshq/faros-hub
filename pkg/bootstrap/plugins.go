package bootstrap

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	apisv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/faroshq/faros-hub/pkg/models"
	utilplugins "github.com/faroshq/faros-hub/pkg/util/plugins"
)

func (b *bootstrap) LoadPlugins(ctx context.Context, workspace string) (models.PluginsList, error) {
	path := b.config.PluginsDir
	var plugins models.PluginsList

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	cluster := logicalcluster.New(b.config.ControllersPluginsWorkspace)

	for _, file := range files {
		p, err := utilplugins.Load(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		name, err := p.GetName(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugin name: %w", err)
		}

		version, err := p.GetVersion(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugin version: %w", err)
		}

		// Register APIResourceSchema for plugin
		data, err := p.GetAPIResourceSchema(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get API resource schema for plugin %q: %w", name, err)
		}

		var schema apisv1alpha1.APIResourceSchema
		err = yaml.Unmarshal(data, &schema)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal API resource schema for plugin %q: %w", name, err)
		}

		_, err = b.kcpClient.Cluster(cluster).ApisV1alpha1().APIResourceSchemas().Create(ctx, &schema, metav1.CreateOptions{})
		if err != nil && apierrors.IsConflict(err) {
			return nil, fmt.Errorf("failed to create API resource schema for plugin %q: %w", name, err)
		}

		// Register APIExport for plugin
		data, err = p.GetAPIExportSchema(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get API export schema for plugin %q: %w", name, err)
		}

		var export apisv1alpha1.APIExport
		err = yaml.Unmarshal(data, &export)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal API export schema for plugin %q: %w", name, err)
		}

		_, err = b.kcpClient.Cluster(cluster).ApisV1alpha1().APIExports().Create(ctx, &export, metav1.CreateOptions{})
		if err != nil && apierrors.IsConflict(err) {
			return nil, fmt.Errorf("failed to create API export schema for plugin %q: %w", name, err)
		}
		plugins = append(plugins, models.Plugin{
			Name:    name,
			Version: version,
		})
	}

	return plugins, nil
}
