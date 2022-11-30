package users

import (
	"context"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

type userDeleteReconciler struct{}

func (r *userDeleteReconciler) reconcile(ctx context.Context, user *tenancyv1alpha1.User) (reconcileStatus, error) {
	return reconcileStatusContinue, nil
}
