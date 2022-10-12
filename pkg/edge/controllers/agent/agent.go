package agent

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
)

// Reconciler reconciles a Potato object
type Reconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Config      *config.Config
	CoreClients kubernetes.ClusterInterface
}

// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent/finalizers,verbs=update

// Reconcile reconciles a Edge object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Include the clusterName from req.ObjectKey in the logger, similar to the namespace and name keys that are already
	// there.
	logger = logger.WithValues("namespace", req.Namespace).WithValues("name", req.Name)

	logger.Info("Getting Agent")
	var agent edgev1alpha1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		if errors.IsNotFound(err) {
			// Normal - was deleted
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	agentCopy := agent.DeepCopy()

	spew.Dump(agentCopy)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&edgev1alpha1.Agent{}).
		Complete(r)
}
