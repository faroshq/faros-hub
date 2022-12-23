package registration

import (
	"context"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type finalizerAddReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerAddReconciler) reconcile(ctx context.Context, _ logicalcluster.Path, registration *edgev1alpha1.Registration) (reconcileStatus, error) {
	if !controllerutil.ContainsFinalizer(registration, r.getFinalizerName()) {
		controllerutil.AddFinalizer(registration, r.getFinalizerName())
		return reconcileStatusStopAndRequeue, nil
	} else {
		return reconcileStatusContinue, nil
	}
}

type finalizerRemoveReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerRemoveReconciler) reconcile(ctx context.Context, _ logicalcluster.Path, registration *edgev1alpha1.Registration) (reconcileStatus, error) {
	controllerutil.RemoveFinalizer(registration, r.getFinalizerName())
	return reconcileStatusContinue, nil
}
