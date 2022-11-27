package users

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/kcp-dev/logicalcluster/v2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

func (r *Reconciler) delete(ctx context.Context, logger logr.Logger, user *tenancyv1alpha1.User, cluster logicalcluster.Name) (ctrl.Result, error) {
	ctx = logicalcluster.WithCluster(ctx, cluster)
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: user.Name,
		},
	}

	err := r.CoreClients.CoreV1().Namespaces().Delete(ctx, namespace.Name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	controllerutil.RemoveFinalizer(user, finalizerName)
	if err := r.Update(ctx, user); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
