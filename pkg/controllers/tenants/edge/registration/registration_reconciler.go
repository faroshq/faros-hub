package registration

import (
	"context"
	"fmt"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/kubernetes"
)

type reconcileStatus int

const (
	reconcileStatusStopAndRequeue reconcileStatus = iota
	reconcileStatusContinue
	reconcileStatusError
)

type reconciler interface {
	reconcile(ctx context.Context, cluster logicalcluster.Name, registration *edgev1alpha1.Registration) (reconcileStatus, error)
}

func (c *Controller) reconcile(ctx context.Context, cluster logicalcluster.Name, registration *edgev1alpha1.Registration) (bool, error) {
	var reconcilers []reconciler
	createReconcilers := []reconciler{
		&finalizerAddReconciler{ // must be first
			getFinalizerName: func() string {
				return finalizerName
			},
		},
		&registrationCreateReconciler{
			getRegistrationResourceName: func(r *edgev1alpha1.Registration) string {
				return getRegistrationResourceName(r)
			},
			createOrUpdateServiceAccount: func(ctx context.Context, cluster logicalcluster.Name, sa *corev1.ServiceAccount) error {
				return createOrUpdateServiceAccount(ctx, c.coreClientSet.Cluster(cluster), sa)
			},
			createOrUpdateSecret: func(ctx context.Context, cluster logicalcluster.Name, secret *corev1.Secret) error {
				return createOrUpdateSecret(ctx, c.coreClientSet.Cluster(cluster), secret)
			},
			createOrUpdateRole: func(ctx context.Context, cluster logicalcluster.Name, role *rbacv1.Role) error {
				return createOrUpdateRole(ctx, c.coreClientSet.Cluster(cluster), role)
			},
			createOrUpdateRoleBinding: func(ctx context.Context, cluster logicalcluster.Name, roleBinding *rbacv1.RoleBinding) error {
				return createOrUpdateRoleBinding(ctx, c.coreClientSet.Cluster(cluster), roleBinding)
			},
			getSecret: func(ctx context.Context, cluster logicalcluster.Name, name, namespace string) (*corev1.Secret, error) {
				return c.coreClientSet.Cluster(cluster).CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
			},
		},
	}

	deleteReconcilers := []reconciler{
		&registrationDeleteReconciler{
			getRegistrationResourceName: func(r *edgev1alpha1.Registration) string {
				return getRegistrationResourceName(r)
			},
			deleteServiceAccount: func(ctx context.Context, cluster logicalcluster.Name, sa *corev1.ServiceAccount) error {
				return c.coreClientSet.Cluster(cluster).CoreV1().ServiceAccounts(sa.Namespace).Delete(ctx, sa.Name, metav1.DeleteOptions{})
			},
			deleteSecret: func(ctx context.Context, cluster logicalcluster.Name, secret *corev1.Secret) error {
				return c.coreClientSet.Cluster(cluster).CoreV1().Secrets(secret.Namespace).Delete(ctx, secret.Name, metav1.DeleteOptions{})
			},
			deleteRole: func(ctx context.Context, cluster logicalcluster.Name, role *rbacv1.Role) error {
				return c.coreClientSet.Cluster(cluster).RbacV1().Roles(role.Namespace).Delete(ctx, role.Name, metav1.DeleteOptions{})
			},
			deleteRoleBinding: func(ctx context.Context, cluster logicalcluster.Name, roleBinding *rbacv1.RoleBinding) error {
				return c.coreClientSet.Cluster(cluster).RbacV1().RoleBindings(roleBinding.Namespace).Delete(ctx, roleBinding.Name, metav1.DeleteOptions{})
			},
		},
		&finalizerRemoveReconciler{
			getFinalizerName: func() string {
				return finalizerName
			},
		},
	}

	if !registration.DeletionTimestamp.IsZero() { //delete
		reconcilers = deleteReconcilers
	} else { //create or update
		reconcilers = createReconcilers
	}

	var errs []error

	requeue := false
	for _, r := range reconcilers {
		var err error
		var status reconcileStatus
		status, err = r.reconcile(ctx, cluster, registration)
		if err != nil {
			errs = append(errs, err)
		}
		if status == reconcileStatusStopAndRequeue {
			requeue = true
			break
		}
	}

	return requeue, utilerrors.NewAggregate(errs)
}

func createOrUpdateServiceAccount(ctx context.Context, coreClients kubernetes.Interface, serviceAccount *corev1.ServiceAccount) error {
	originalOwners := serviceAccount.OwnerReferences

	currentServiceAccount, err := coreClients.CoreV1().ServiceAccounts(serviceAccount.Namespace).Get(ctx, serviceAccount.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err := coreClients.CoreV1().ServiceAccounts(serviceAccount.Namespace).Create(ctx, serviceAccount, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create the ServiceAccount %s", err)
		}
	case err == nil:
		currentServiceAccount.ResourceVersion = ""
		currentServiceAccount.OwnerReferences = mergeOwnerReference(currentServiceAccount.OwnerReferences, originalOwners)
		_, err := coreClients.CoreV1().ServiceAccounts(serviceAccount.Namespace).Update(ctx, currentServiceAccount, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update the ServiceAccount %s", err)
		}
	default:
		return fmt.Errorf("failed to create the ServiceAccount %s", err)
	}

	return nil
}

func createOrUpdateSecret(ctx context.Context, coreClients kubernetes.Interface, secret *corev1.Secret) error {
	originalOwners := secret.OwnerReferences

	currentSecret, err := coreClients.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err := coreClients.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create the ServiceAccount %s", err)
		}
	case err == nil:
		currentSecret.ResourceVersion = ""
		currentSecret.OwnerReferences = mergeOwnerReference(currentSecret.OwnerReferences, originalOwners)
		_, err := coreClients.CoreV1().Secrets(currentSecret.Namespace).Update(ctx, currentSecret, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update the Secret %s", err)
		}
	default:
		return fmt.Errorf("failed to create the Secret %s", err)
	}

	return nil
}

func createOrUpdateRole(ctx context.Context, coreClients kubernetes.Interface, role *rbacv1.Role) error {
	originalOwners := role.OwnerReferences

	currentRole, err := coreClients.RbacV1().Roles(role.Namespace).Get(ctx, role.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err := coreClients.RbacV1().Roles(role.Namespace).Create(ctx, role, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create the Role %s", err)
		}
	case err == nil:
		currentRole.Rules = role.Rules
		currentRole.ResourceVersion = ""
		currentRole.OwnerReferences = mergeOwnerReference(currentRole.OwnerReferences, originalOwners)
		_, err := coreClients.RbacV1().Roles(role.Namespace).Update(ctx, currentRole, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update the Role %s", err)
		}
	default:
		return fmt.Errorf("failed to create the Role %s", err)
	}

	return nil
}

func createOrUpdateRoleBinding(ctx context.Context, coreClients kubernetes.Interface, roleBinding *rbacv1.RoleBinding) error {
	originalOwners := roleBinding.OwnerReferences

	currentRoleBinding, err := coreClients.RbacV1().RoleBindings(roleBinding.Namespace).Get(ctx, roleBinding.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err := coreClients.RbacV1().RoleBindings(roleBinding.Namespace).Create(ctx, roleBinding, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create the RoleBindings %s", err)
		}
	case err == nil:
		currentRoleBinding.RoleRef = roleBinding.RoleRef
		currentRoleBinding.Subjects = roleBinding.Subjects
		currentRoleBinding.ResourceVersion = ""
		currentRoleBinding.OwnerReferences = mergeOwnerReference(currentRoleBinding.OwnerReferences, originalOwners)
		_, err := coreClients.RbacV1().RoleBindings(roleBinding.Namespace).Update(ctx, currentRoleBinding, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update the RoleBindings %s", err)
		}
	default:
		return fmt.Errorf("failed to create the RoleBindings %s", err)
	}

	return nil
}
