package registration

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/kcp-dev/logicalcluster/v3"
)

type registrationCreateReconciler struct {
	getRegistrationResourceName  func(r *edgev1alpha1.Registration) string
	createOrUpdateServiceAccount func(ctx context.Context, cluster logicalcluster.Path, sa *corev1.ServiceAccount) error
	createOrUpdateSecret         func(ctx context.Context, cluster logicalcluster.Path, secret *corev1.Secret) error
	createOrUpdateRole           func(ctx context.Context, cluster logicalcluster.Path, role *rbacv1.Role) error
	createOrUpdateRoleBinding    func(ctx context.Context, cluster logicalcluster.Path, roleBinding *rbacv1.RoleBinding) error
	getSecret                    func(ctx context.Context, cluster logicalcluster.Path, name, namespace string) (*corev1.Secret, error)
}

func (r *registrationCreateReconciler) reconcile(ctx context.Context, cluster logicalcluster.Path, registration *edgev1alpha1.Registration) (reconcileStatus, error) {
	resourceName := r.getRegistrationResourceName(registration)
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	}

	err := r.createOrUpdateServiceAccount(ctx, cluster, sa)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// Dedicates secret name to the registration
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": resourceName,
			},
			OwnerReferences: getRegistrationOwnersReference(registration),
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}

	err = r.createOrUpdateSecret(ctx, cluster, secret)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// Create a cluster role that provides the agent the minimal permissions
	// required by KCP to manage the agent target, and by the agent
	rules := []rbacv1.PolicyRule{
		{
			Verbs:     []string{"create", "list", "get", "watch"},
			APIGroups: []string{edgev1alpha1.SchemeGroupVersion.Group},
			Resources: []string{"agents"},
		},
		{
			Verbs:     []string{"update", "patch"},
			APIGroups: []string{edgev1alpha1.SchemeGroupVersion.Group},
			Resources: []string{"agents/status"},
		},
	}

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:            resourceName,
			Namespace:       registration.Namespace,
			OwnerReferences: getRegistrationOwnersReference(registration),
		},
		Rules: rules,
	}

	err = r.createOrUpdateRole(ctx, cluster, role)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// Add role binding

	// Grant the service account the role created just above in the workspace
	subjects := []rbacv1.Subject{{
		Kind:      "ServiceAccount",
		Name:      resourceName,
		Namespace: registration.Namespace,
	}}
	roleRef := rbacv1.RoleRef{
		Kind:     "Role",
		Name:     resourceName,
		APIGroup: "rbac.authorization.k8s.io",
	}

	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:            resourceName,
			Namespace:       registration.Namespace,
			OwnerReferences: getRegistrationOwnersReference(registration),
		},
		Subjects: subjects,
		RoleRef:  roleRef,
	}

	err = r.createOrUpdateRoleBinding(ctx, cluster, roleBinding)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// Retrieve the token that the agent will use to authenticate to kcp
	tokenSecret, err := r.getSecret(ctx, cluster, resourceName, registration.Namespace)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	saTokenBytes, ok := tokenSecret.Data["token"]
	if !ok {
		return reconcileStatusStopAndRequeue, fmt.Errorf("token not found in secret %s", resourceName)
	}
	if len(saTokenBytes) == 0 {
		return reconcileStatusStopAndRequeue, fmt.Errorf("token secret %s/%s is missing a value for `token`", registration.Namespace, resourceName)
	}

	registration.Status.Token = string(saTokenBytes)
	conditions.MarkTrue(registration, conditionsv1alpha1.ReadyCondition)

	return reconcileStatusContinue, nil
}
