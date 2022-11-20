package registration

import (
	"context"
	"fmt"
	"time"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/go-logr/logr"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	workloadv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/workload/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Reconciler) createOrUpdate(ctx context.Context, logger logr.Logger, registration *edgev1alpha1.Registration) (ctrl.Result, error) {
	// TODO: move to webhook
	if !controllerutil.ContainsFinalizer(registration, finalizerName) {
		controllerutil.AddFinalizer(registration, finalizerName)
		if err := r.Update(ctx, registration); err != nil {
			return ctrl.Result{}, err
		}
		// requeue to ensure the finalizer is set before creating the resources
		return ctrl.Result{Requeue: true}, nil
	}

	registrationOwnersReferences := []metav1.OwnerReference{{
		APIVersion: workloadv1alpha1.SchemeGroupVersion.String(),
		Kind:       edgev1alpha1.RegistrationKind,
		Name:       registration.Name,
		UID:        registration.UID,
	}}

	// ServiceAccount creation
	resourceName := getRegistrationResourceName(registration.Name)
	var sa corev1.ServiceAccount
	err := r.Get(ctx, client.ObjectKey{
		Namespace: registration.Namespace,
		Name:      resourceName,
	}, &sa)
	switch {
	case apierrors.IsNotFound(err):
		logger.Error(err, "creating service account", "service-account-name", resourceName)
		err := r.Create(ctx, &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:            resourceName,
				Namespace:       registration.Namespace,
				OwnerReferences: registrationOwnersReferences,
			},
		}, &client.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create ServiceAccount: %s", err)
		}
	case err == nil:
		// service account already exist, merge owner references
		sa.OwnerReferences = mergeOwnerReference(sa.OwnerReferences, registrationOwnersReferences)
		sa.ResourceVersion = ""

		logger.Info("updating service account", "service-account-name", resourceName)
		err = r.Patch(ctx, &sa, client.MergeFrom(sa.DeepCopy()), &client.PatchOptions{})
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to patch ServiceAccount %s", err)
		}
	default:
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to get the ServiceAccount %s", err)
	}

	// Dedicates secret name to the registration
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": resourceName,
			},
			OwnerReferences: registrationOwnersReferences,
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}
	err = r.Get(ctx, client.ObjectKey{
		Namespace: secret.Namespace,
		Name:      resourceName,
	}, secret)
	switch {
	case apierrors.IsNotFound(err):
		logger.Error(err, "creating secret", "name", resourceName)
		err := r.Create(ctx, secret, &client.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create Secret: %s", err)
		}
	case err == nil:
		// service account already exist, merge owner references
		secret.OwnerReferences = mergeOwnerReference(secret.OwnerReferences, registrationOwnersReferences)
		secret.ResourceVersion = ""

		logger.Info("updating secret", "name", resourceName)
		err = r.Patch(ctx, secret, client.MergeFrom(secret.DeepCopy()), &client.PatchOptions{})
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to patch Secret %s", err)
		}
	default:
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to get the Secret %s", err)
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

	role := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:            resourceName,
			Namespace:       registration.Namespace,
			OwnerReferences: registrationOwnersReferences,
		},
		Rules: rules,
	}

	err = r.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: registration.Namespace}, &role)
	switch {
	case apierrors.IsNotFound(err):
		logger.Info("creating role", "role-name", role.Name)
		err := r.Create(ctx, &role, &client.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create the Role %s", err)
		}
	case err == nil:
		role.Rules = rules
		role.ResourceVersion = ""
		role.OwnerReferences = mergeOwnerReference(role.OwnerReferences, registrationOwnersReferences)
		err := r.Update(ctx, &role, &client.UpdateOptions{})
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to update the Role %s", err)
		}
	default:
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create the Role %s", err)
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

	roleBinding := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:            resourceName,
			Namespace:       registration.Namespace,
			OwnerReferences: registrationOwnersReferences,
		},
		Subjects: subjects,
		RoleRef:  roleRef,
	}
	err = r.Get(ctx, types.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}, &roleBinding)
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to get RoleBinding %s", err)
	}
	if err == nil {
		logger.Info("cleaning old role binding", "rolebinding-name", resourceName)
		err := r.Delete(ctx, &roleBinding, &client.DeleteOptions{})
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to delete RoleBinding %s", err)
		}
	}

	logger.Info("creating or updating role binding to bind service account to cluster role", "rolebinding-name", resourceName)
	roleBinding.ResourceVersion = ""
	err = r.Create(ctx, &roleBinding, &client.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create/update RoleBinding %s", err)
	}

	// Retrieve the token that the agent will use to authenticate to kcp
	tokenSecret := corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: registration.Namespace}, &tokenSecret)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to retrieve Secret: %w", err)
	}
	saTokenBytes, ok := tokenSecret.Data["token"]
	if !ok {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to retrieve token from secret: %w", err)
	}
	if len(saTokenBytes) == 0 {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("token secret %s/%s is missing a value for `token`", registration.Namespace, resourceName)
	}
	//caCrt, ok := tokenSecret.Data["ca.crt"]
	//////if !ok {
	//////	return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to retrieve ca.crt from secret: %w", err)
	//}

	patch := client.MergeFrom(registration.DeepCopy())
	registration.Status.Token = string(saTokenBytes)
	conditions.MarkTrue(registration, conditionsv1alpha1.ReadyCondition)

	if err := r.Status().Patch(ctx, registration, patch); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
