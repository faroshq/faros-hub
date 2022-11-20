package registration

import (
	"context"
	"fmt"
	"time"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Reconciler) delete(ctx context.Context, logger logr.Logger, registration *edgev1alpha1.Registration) (ctrl.Result, error) {

	resourceName := getRegistrationResourceName(registration.Name)
	// ServiceAccount deletion
	err := r.Delete(ctx, &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	})
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to delete ServiceAccount: %s", err)
	}

	// Secret deletion
	err = r.Delete(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	})
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to delete Secret: %s", err)
	}

	// Role deletion
	err = r.Delete(ctx, &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	})
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to delete Role: %s", err)
	}

	// RoleBinding deletion
	err = r.Delete(ctx, &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: registration.Namespace,
		},
	})
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to delete RoleBinding: %s", err)
	}

	controllerutil.RemoveFinalizer(registration, finalizerName)
	if err := r.Update(ctx, registration); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
