package pluginbindings

import (
	"context"
	"fmt"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
)

// Binding needs to:
// Know if plugin is installed in the cluster (check if APIExport exists)
// If it does we need to select all agents in the namespace and update agents status.

type bindingReconciler struct {
	getAPIBinding func(ctx context.Context, cluster logicalcluster.Path, name string) (bool, error)
	//getAgents     func(ctx context.Context, cluster logicalcluster.Name, namespace string, selector metav1.LabelSelector) ([]string, error)
}

func (r *bindingReconciler) reconcile(ctx context.Context, cluster logicalcluster.Path, binding *pluginsv1alpha1.Binding) (reconcileStatus, error) {
	// Check if APIExport exists
	apiExportExists, err := r.getAPIBinding(ctx, cluster, binding.Name)
	if err != nil {
		return reconcileStatusError, err
	}
	if !apiExportExists {
		return reconcileStatusStopAndRequeue, fmt.Errorf("plugin %s/%s is not enabled", binding.Spec.PluginType, binding.Spec.PluginName)
	}

	// Select all agents in the namespace and update agents status.
	//getAgents(ctx)

	return reconcileStatusContinue, nil
}
