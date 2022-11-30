package users

import (
	"context"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
)

type userCreateReconciler struct{}

func (r *userCreateReconciler) reconcile(ctx context.Context, user *tenancyv1alpha1.User) (reconcileStatus, error) {
	conditions.MarkTrue(user, conditionsv1alpha1.ReadyCondition)
	return reconcileStatusContinue, nil
}
