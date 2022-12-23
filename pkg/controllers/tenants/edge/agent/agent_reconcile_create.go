package agent

import (
	"context"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
)

type agentCreateReconciler struct {
}

func (r *agentCreateReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, agent *edgev1alpha1.Agent) (reconcileStatus, error) {

	return reconcileStatusContinue, nil
}
