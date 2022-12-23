package pluginbindings

import (
	"context"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	kcpclientset "github.com/kcp-dev/kcp/pkg/client/clientset/versioned/cluster"
	"github.com/kcp-dev/logicalcluster/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type reconcileStatus int

const (
	reconcileStatusStopAndRequeue reconcileStatus = iota
	reconcileStatusContinue
	reconcileStatusError
)

type reconciler interface {
	reconcile(ctx context.Context, cluster logicalcluster.Path, binding *pluginsv1alpha1.Binding) (reconcileStatus, error)
}

func (c *Controller) reconcile(ctx context.Context, cluster logicalcluster.Path, binding *pluginsv1alpha1.Binding) (bool, error) {
	var reconcilers []reconciler
	createReconcilers := []reconciler{
		&finalizerAddReconciler{ // must be first
			getFinalizerName: func() string {
				return finalizerName
			},
		},
		&bindingReconciler{
			getAPIBinding: func(ctx context.Context, cluster logicalcluster.Path, name string) (bool, error) {
				return getAPIBinding(ctx, c.kcpClientSet, cluster, name)
			},
		},
	}

	deleteReconcilers := []reconciler{
		&finalizerRemoveReconciler{
			getFinalizerName: func() string {
				return finalizerName
			},
		},
		&bindingReconciler{},
	}

	if !binding.DeletionTimestamp.IsZero() { //delete
		reconcilers = deleteReconcilers
	} else { //create or update
		reconcilers = createReconcilers
	}

	var errs []error

	requeue := false
	for _, r := range reconcilers {
		var err error
		var status reconcileStatus
		status, err = r.reconcile(ctx, cluster, binding)
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

func getAPIBinding(ctx context.Context, kcpClient kcpclientset.ClusterInterface, cluster logicalcluster.Path, name string) (bool, error) {
	_, err := kcpClient.Cluster(cluster).ApisV1alpha1().APIExports().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return false, err
	}
	return !apierrors.IsNotFound(err), nil
}
