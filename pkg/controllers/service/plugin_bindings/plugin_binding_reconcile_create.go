package pluginbindings

import (
	"context"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
)

type bindingReconciler struct {
}

func (r *bindingReconciler) reconcile(ctx context.Context, cluster logicalcluster.Name, binding *pluginsv1alpha1.Binding) (reconcileStatus, error) {

	return reconcileStatusContinue, nil
}
