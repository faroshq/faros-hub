package users

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

type userDeleteReconciler struct {
	deleteNamespace func(ctx context.Context, namespace *corev1.Namespace) error
}

func (r *userDeleteReconciler) reconcile(ctx context.Context, user *tenancyv1alpha1.User) (reconcileStatus, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:            user.Name,
			OwnerReferences: getWorkspaceOwnersReference(user),
		},
	}

	err := r.deleteNamespace(ctx, namespace)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	return reconcileStatusContinue, nil
}
