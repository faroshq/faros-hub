package agent

import (
	"context"

	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/kcp-dev/logicalcluster/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/models"
)

var finalizerName = "agents.edge.faros.sh/finalizer"

// Reconciler reconciles a Agent object
type Reconciler struct {
	client.Client
	RootRest *rest.Config
	Scheme   *runtime.Scheme
	Plugins  models.PluginsList
}

// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent/finalizers,verbs=update

// Reconcile reconciles a Agent object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Include the clusterName from req.ObjectKey in the logger, similar to the namespace and name keys that are already
	// there.
	logger = logger.WithValues("clusterName", req.ClusterName).WithValues("namespace", req.Namespace).WithValues("name", req.Name)

	// Add the logical cluster to the context
	ctx = logicalcluster.WithCluster(ctx, logicalcluster.New(req.ClusterName))

	var agent edgev1alpha1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var result ctrl.Result
	var err error
	if agent.DeletionTimestamp.IsZero() {
		result, err = r.createOrUpdate(ctx, logger, agent.DeepCopy())
	} else {
		result, err = r.delete(ctx, logger, agent.DeepCopy())
	}
	if err != nil {
		agentCopy := agent.DeepCopy()
		conditions.MarkFalse(
			agentCopy,
			conditionsv1alpha1.ReadyCondition,
			err.Error(),
			conditionsv1alpha1.ConditionSeverityError,
			"Error configuring Agent: %v",
			err,
		)
		if err := r.Status().Patch(ctx, agentCopy, client.MergeFrom(&agent)); err != nil {
			return result, err
		}
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&edgev1alpha1.Agent{}).
		Complete(r)
}
