package agent

import (
	"fmt"
	"path/filepath"
	"runtime"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/plugins"

	utilplugins "github.com/faroshq/faros-hub/pkg/util/plugins"
)

func (r *Reconciler) loadPlugin(plugin edgev1alpha1.PluginSpec) (plugins.Interface, error) {
	path := filepath.Join(r.Config.PluginsDir, pluginFileName(plugin))
	return utilplugins.Load(path)
}

func pluginFileName(plugin edgev1alpha1.PluginSpec) string {
	return fmt.Sprintf("%s-%s-%s", plugin.Name, plugin.Version, runtime.GOARCH)
}
