package pluginrequests

import (
	"context"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
)

type requestDeleteReconciler struct{}

func (r *requestDeleteReconciler) reconcile(ctx context.Context, cluster logicalcluster.Path, request *pluginsv1alpha1.Request) (reconcileStatus, error) {
	return reconcileStatusContinue, nil
}
