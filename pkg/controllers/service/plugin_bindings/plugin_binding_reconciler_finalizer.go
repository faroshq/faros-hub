package pluginbindings

import (
	"context"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type finalizerAddReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerAddReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, binding *pluginsv1alpha1.Binding) (reconcileStatus, error) {
	if !controllerutil.ContainsFinalizer(binding, r.getFinalizerName()) {
		controllerutil.AddFinalizer(binding, r.getFinalizerName())
		return reconcileStatusStopAndRequeue, nil
	} else {
		return reconcileStatusContinue, nil
	}
}

type finalizerRemoveReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerRemoveReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, binding *pluginsv1alpha1.Binding) (reconcileStatus, error) {
	controllerutil.RemoveFinalizer(binding, r.getFinalizerName())
	return reconcileStatusContinue, nil
}
