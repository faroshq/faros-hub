package agent

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/go-logr/logr"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Reconciler) createOrUpdate(ctx context.Context, logger logr.Logger, agent *edgev1alpha1.Agent) (ctrl.Result, error) {
	// TODO: move to webhook
	if !controllerutil.ContainsFinalizer(agent, finalizerName) {
		controllerutil.AddFinalizer(agent, finalizerName)
		if err := r.Update(ctx, agent); err != nil {
			return ctrl.Result{}, err
		}
		// requeue to ensure the finalizer is set before creating the resources
		return ctrl.Result{Requeue: true}, nil
	}

	// Iterate plugins and see if they are available and create instances of them
	for _, plugin := range agent.Spec.Plugins {

		if plugin.Version == "latest" {
			_, err := r.Plugins.GetLatest(plugin.Name)
			if err != nil {
				logger.Error(err, "Failed to get latest version of plugin", "plugin", plugin.Name)
				conditions.MarkFalse(agent, conditionsv1alpha1.ReadyCondition, "PluginNotFound", "Failed to get latest version of plugin %s", plugin.Name)
				return ctrl.Result{}, err
			}
			// Create a registration for each plugin
		} else {
			if r.Plugins.Has(models.Plugin{
				Name:    plugin.Name,
				Version: plugin.Version,
			}) {
				spew.Dump(models.Plugin{
					Name:    plugin.Name,
					Version: plugin.Version,
				})
			} else {
				// Mark the plugin as unavailable
				conditions.MarkFalse(agent, conditionsv1alpha1.ReadyCondition, "PluginUnavailable", "Plugin %s is not available", plugin.Name)
				if err := r.Update(ctx, agent); err != nil {
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, nil
			}
		}
	}

	patch := client.MergeFrom(agent.DeepCopy())
	conditions.MarkTrue(agent, conditionsv1alpha1.ReadyCondition)

	if err := r.Status().Patch(ctx, agent, patch); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
