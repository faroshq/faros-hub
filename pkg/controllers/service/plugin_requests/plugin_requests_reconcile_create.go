package pluginrequests

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/models"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/kcp-dev/logicalcluster/v2"
)

type requestCreateReconciler struct {
	getPlugins      func() models.PluginsList
	createAPIExport func(ctx context.Context, destinationCluster logicalcluster.Name, pluginVersion, pluginName string) error
}

func (r *requestCreateReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, request *pluginsv1alpha1.Request) (reconcileStatus, error) {
	conditions.MarkTrue(request, conditionsv1alpha1.ReadyCondition)

	availablePlugins := r.getPlugins()

	spew.Dump(request)

	pluginName := request.Spec.Name
	pluginVersion := request.Spec.Version

	if pluginVersion == "latest" {
		plugin, err := availablePlugins.GetLatest(pluginName)
		if err != nil {
			return reconcileStatusError, err
		}
		pluginVersion = plugin.Version
	} else {
		exists := availablePlugins.Has(models.Plugin{
			Name:    pluginName,
			Version: pluginVersion,
		})
		if !exists {
			return reconcileStatusError, fmt.Errorf("plugin %s:%s does not exist", pluginName, pluginVersion)
		}
	}

	// enable plugin/apibinding in the the cluster
	err := r.createAPIExport(ctx, cluster, pluginVersion, pluginName)
	if err != nil {
		return reconcileStatusError, err
	}

	request.Status.Version = pluginVersion

	return reconcileStatusContinue, nil
}
