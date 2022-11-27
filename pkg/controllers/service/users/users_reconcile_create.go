package users

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
)

type userCreateReconciler struct {
	createNamespace func(ctx context.Context, namespace *corev1.Namespace) error
}

func (r *userCreateReconciler) reconcile(ctx context.Context, user *tenancyv1alpha1.User) (reconcileStatus, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:            user.Name,
			OwnerReferences: getWorkspaceOwnersReference(user),
		},
	}

	err := r.createNamespace(ctx, namespace)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	conditions.MarkTrue(user, conditionsv1alpha1.ReadyCondition)

	return reconcileStatusContinue, nil
}
