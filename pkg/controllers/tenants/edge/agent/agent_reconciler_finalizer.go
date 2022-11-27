package agent

import (
	"context"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type finalizerAddReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerAddReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, agent *edgev1alpha1.Agent) (reconcileStatus, error) {
	if !controllerutil.ContainsFinalizer(agent, r.getFinalizerName()) {
		controllerutil.AddFinalizer(agent, r.getFinalizerName())
		return reconcileStatusStopAndRequeue, nil
	} else {
		return reconcileStatusContinue, nil
	}
}

type finalizerRemoveReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerRemoveReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, agent *edgev1alpha1.Agent) (reconcileStatus, error) {
	controllerutil.RemoveFinalizer(agent, r.getFinalizerName())
	return reconcileStatusContinue, nil
}
