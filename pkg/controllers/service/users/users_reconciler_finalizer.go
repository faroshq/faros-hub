package users

import (
	"context"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type finalizerAddReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerAddReconciler) reconcile(ctx context.Context, user *tenancyv1alpha1.User) (reconcileStatus, error) {
	if !controllerutil.ContainsFinalizer(user, r.getFinalizerName()) {
		controllerutil.AddFinalizer(user, r.getFinalizerName())
		return reconcileStatusStopAndRequeue, nil
	} else {
		return reconcileStatusContinue, nil
	}
}

type finalizerRemoveReconciler struct {
	getFinalizerName func() string
}

func (r *finalizerRemoveReconciler) reconcile(ctx context.Context, user *tenancyv1alpha1.User) (reconcileStatus, error) {
	controllerutil.RemoveFinalizer(user, r.getFinalizerName())
	return reconcileStatusContinue, nil
}
