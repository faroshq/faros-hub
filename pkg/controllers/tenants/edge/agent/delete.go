package agent

import (
	"context"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Reconciler) delete(ctx context.Context, logger logr.Logger, agent *edgev1alpha1.Agent) (ctrl.Result, error) {

	return ctrl.Result{}, nil
}
