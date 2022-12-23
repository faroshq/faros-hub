package agent

import (
	"context"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type reconcileStatus int

const (
	reconcileStatusStopAndRequeue reconcileStatus = iota
	reconcileStatusContinue
	reconcileStatusError
)

type reconciler interface {
	reconcile(ctx context.Context, cluster logicalcluster.Name, agent *edgev1alpha1.Agent) (reconcileStatus, error)
}

func (c *Controller) reconcile(ctx context.Context, cluster logicalcluster.Name, agent *edgev1alpha1.Agent) (bool, error) {
	var reconcilers []reconciler
	createReconcilers := []reconciler{
		&finalizerAddReconciler{ // must be first
			getFinalizerName: func() string {
				return finalizerName
			},
		},
		&agentCreateReconciler{},
	}

	deleteReconcilers := []reconciler{
		&agentCreateReconciler{},
		&finalizerRemoveReconciler{
			getFinalizerName: func() string {
				return finalizerName
			},
		},
	}

	if !agent.DeletionTimestamp.IsZero() { //delete
		reconcilers = deleteReconcilers
	} else { //create or update
		reconcilers = createReconcilers
	}

	var errs []error

	requeue := false
	for _, r := range reconcilers {
		var err error
		var status reconcileStatus
		status, err = r.reconcile(ctx, cluster, agent)
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
