package bootstrap

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/davecgh/go-spew/spew"
	apisv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"

	"github.com/faroshq/faros-hub/pkg/models"
	utilplugins "github.com/faroshq/faros-hub/pkg/util/plugins"
)

type pluginStore struct {
	plugins models.PluginsList
	lock    sync.RWMutex
}

func (s *pluginStore) Get() models.PluginsList {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.plugins
}

func (s *pluginStore) Set(plugin models.Plugin) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.plugins = append(s.plugins, plugin)
}

func newPluginStore() *pluginStore {
	return &pluginStore{
		plugins: models.PluginsList{},
		lock:    sync.RWMutex{},
	}
}

func (b *bootstrap) LoadPlugins(ctx context.Context, workspace string) (models.PluginsList, error) {
	logger := klog.FromContext(ctx)

	clusterPath := logicalcluster.NewPath(workspace)
	path := b.config.PluginsDir
	store := newPluginStore()

	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(err, "failed to read plugins directory. No plugins will be served")
		return nil, nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(files))

	for _, file := range files {
		f := func(file fs.FileInfo, wg *sync.WaitGroup) error {
			p, err := utilplugins.Load(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}

			name, err := p.GetName(ctx)
			if err != nil {
				return fmt.Errorf("failed to get plugin name: %w", err)
			}

			version, err := p.GetVersion(ctx)
			if err != nil {
				return fmt.Errorf("failed to get plugin version: %w", err)
			}

			// Register APIResourceSchema for plugin
			data, err := p.GetAPIResourceSchema(ctx)
			if err != nil {
				return fmt.Errorf("failed to get API resource schema for plugin %q: %w", name, err)
			}

			var schema apisv1alpha1.APIResourceSchema
			err = yaml.Unmarshal(data, &schema)
			if err != nil {
				return fmt.Errorf("failed to unmarshal API resource schema for plugin %q: %w", name, err)
			}

			_, err = b.kcpClient.Cluster(clusterPath).ApisV1alpha1().APIResourceSchemas().Create(ctx, &schema, metav1.CreateOptions{})
			if err != nil && apierrors.IsConflict(err) {
				return fmt.Errorf("failed to create API resource schema for plugin %q: %w", name, err)
			}

			// Register APIExport for plugin
			data, err = p.GetAPIExportSchema(ctx)
			if err != nil {
				return fmt.Errorf("failed to get API export schema for plugin %q: %w", name, err)
			}

			var export apisv1alpha1.APIExport
			err = yaml.Unmarshal(data, &export)
			if err != nil {
				return fmt.Errorf("failed to unmarshal API export schema for plugin %q: %w", name, err)
			}
			spew.Dump(export)
			spew.Dump(string(data))

			_, err = b.kcpClient.Cluster(clusterPath).ApisV1alpha1().APIExports().Create(ctx, &export, metav1.CreateOptions{})
			if err != nil && apierrors.IsConflict(err) {
				return fmt.Errorf("failed to create API export schema for plugin %q: %w", name, err)
			}
			store.Set(models.Plugin{
				Name:    name,
				Version: version,
			})
			return nil
		}

		go func(file fs.FileInfo, wg *sync.WaitGroup) {
			defer wg.Done()
			if err := f(file, wg); err != nil {
				klog.Error(err)
			}
		}(file, wg)
	}

	wg.Wait()

	return store.Get(), nil
}
