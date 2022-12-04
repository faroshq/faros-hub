package pluginbindings

import (
	"context"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type reconcileStatus int

const (
	reconcileStatusStopAndRequeue reconcileStatus = iota
	reconcileStatusContinue
	reconcileStatusError
)

type reconciler interface {
	reconcile(ctx context.Context, cluster logicalcluster.Name, binding *pluginsv1alpha1.Binding) (reconcileStatus, error)
}

func (c *Controller) reconcile(ctx context.Context, cluster logicalcluster.Name, binding *pluginsv1alpha1.Binding) (bool, error) {
	var reconcilers []reconciler
	createReconcilers := []reconciler{
		&finalizerAddReconciler{ // must be first
			getFinalizerName: func() string {
				return finalizerName
			},
		},
		&bindingReconciler{},
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
