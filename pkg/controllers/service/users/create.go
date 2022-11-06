package users

import (
	"context"

	"github.com/go-logr/logr"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/kcp-dev/logicalcluster/v2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

func (r *Reconciler) createOrUpdate(ctx context.Context, logger logr.Logger, user *tenancyv1alpha1.User, cluster logicalcluster.Name) (ctrl.Result, error) {
	// TODO: move to webhook
	if !controllerutil.ContainsFinalizer(user, finalizerName) {
		controllerutil.AddFinalizer(user, finalizerName)
		if err := r.Update(ctx, user); err != nil {
			return ctrl.Result{}, err
		}
		// requeue to ensure the finalizer is set before creating the resources
		return ctrl.Result{Requeue: true}, nil
	}

	ownersReferences := []metav1.OwnerReference{{
		APIVersion:         tenancyv1alpha1.SchemeGroupVersion.String(),
		Kind:               tenancyv1alpha1.UserKind,
		Name:               user.Name,
		BlockOwnerDeletion: pointer.BoolPtr(true),
		Controller:         pointer.BoolPtr(true),
		UID:                user.UID,
	}}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:            user.Name,
			OwnerReferences: ownersReferences,
		},
	}

	// TODO: Original controller-runtime client discovery is somehow broken. need to debug more
	_, err := r.CoreClients.Cluster(cluster).CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	patch := client.MergeFrom(user.DeepCopy())
	conditions.MarkTrue(user, conditionsv1alpha1.ReadyCondition)

	if err := r.Status().Patch(ctx, user, patch); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
