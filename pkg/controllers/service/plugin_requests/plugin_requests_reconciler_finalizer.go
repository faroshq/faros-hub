package pluginrequests

import (
	"context"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type finalizerAddReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerAddReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, request *pluginsv1alpha1.Request) (reconcileStatus, error) {
	if !controllerutil.ContainsFinalizer(request, r.getFinalizerName()) {
		controllerutil.AddFinalizer(request, r.getFinalizerName())
		return reconcileStatusStopAndRequeue, nil
	} else {
		return reconcileStatusContinue, nil
	}
}

type finalizerRemoveReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerRemoveReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, request *pluginsv1alpha1.Request) (reconcileStatus, error) {
	controllerutil.RemoveFinalizer(request, r.getFinalizerName())
	return reconcileStatusContinue, nil
}
