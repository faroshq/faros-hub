package workspaces

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/faroshq/faros-hub/pkg/config"
	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
)

// Workspaces controller runs on global level and makes sure APIBindings are
// created for all faros workspaces. In the future it will be responsible for
// lifecycle those bindings too.

// Reconciler reconciles an object
type Reconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Config        *config.ControllerConfig
	ClientFactory utilkubernetes.ClientFactory
	CoreClients   kubernetes.ClusterInterface
}

// Reconcile reconciles an object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Include the clusterName from req.ObjectKey in the logger, similar to the namespace and name keys that are already
	// there.
	logger = logger.WithValues("clusterName", req.ClusterName).WithValues("namespace", req.Namespace).WithValues("name", req.Name)

	// Add the logical cluster to the context
	ctx = logicalcluster.WithCluster(ctx, logicalcluster.New(req.ClusterName))

	logger.Info("Getting Request")
	var request tenancyv1alpha1.ClusterWorkspace
	if err := r.Get(ctx, req.NamespacedName, &request); err != nil {
		if errors.IsNotFound(err) {
			// Normal - was deleted
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	requestCopy := request.DeepCopy()

	spew.Dump(requestCopy)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenancyv1alpha1.ClusterWorkspace{}).
		Complete(r)
}
