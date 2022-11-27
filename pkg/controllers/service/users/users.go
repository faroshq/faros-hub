package users

import (
	"context"

	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
)

// Workspaces controller runs on global level and makes sure APIBindings are
// created for all faros workspaces. In the future it will be responsible for
// lifecycle those bindings too.

var finalizerName = "users.tenancy.faros.sh/finalizer"

// Reconciler reconciles an object
type Reconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Config      *config.ControllerConfig
	CoreClients kubernetes.Interface
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
	var request tenancyv1alpha1.User
	if err := r.Get(ctx, req.NamespacedName, &request); err != nil {
		if errors.IsNotFound(err) {
			// Normal - was deleted
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var result ctrl.Result
	var err error
	if request.DeletionTimestamp.IsZero() {
		result, err = r.createOrUpdate(ctx, logger, request.DeepCopy(), logicalcluster.New(req.ClusterName))
	} else {
		result, err = r.delete(ctx, logger, request.DeepCopy(), logicalcluster.New(req.ClusterName))
	}
	if err != nil {
		requestCopy := request.DeepCopy()
		conditions.MarkFalse(
			requestCopy,
			conditionsv1alpha1.ReadyCondition,
			err.Error(),
			conditionsv1alpha1.ConditionSeverityError,
			"Error configuring User: %v",
			err,
		)
		if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
			return result, err
		}
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenancyv1alpha1.User{}).
		Complete(r)
}
