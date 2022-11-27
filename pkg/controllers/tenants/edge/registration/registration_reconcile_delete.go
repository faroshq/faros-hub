package registration

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
)

type registrationDeleteReconciler struct {
	getRegistrationResourceName func(r *edgev1alpha1.Registration) string
	deleteServiceAccount        func(ctx context.Context, cluster logicalcluster.Name, sa *corev1.ServiceAccount) error
	deleteSecret                func(ctx context.Context, cluster logicalcluster.Name, secret *corev1.Secret) error
	deleteRole                  func(ctx context.Context, cluster logicalcluster.Name, role *rbacv1.Role) error
	deleteRoleBinding           func(ctx context.Context, cluster logicalcluster.Name, roleBinding *rbacv1.RoleBinding) error
}

func (r *registrationDeleteReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, registration *edgev1alpha1.Registration) (reconcileStatus, error) {
	resourceName := r.getRegistrationResourceName(registration)
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	}

	err := r.deleteServiceAccount(ctx, cluster, sa)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// Dedicates secret name to the registration
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	}

	err = r.deleteSecret(ctx, cluster, secret)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	}

	err = r.deleteRole(ctx, cluster, role)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	}

	err = r.deleteRoleBinding(ctx, cluster, roleBinding)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	return reconcileStatusContinue, nil
}
