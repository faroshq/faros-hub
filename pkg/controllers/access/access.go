package access

import (
	"context"

	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	accessv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/access/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
	kubernetesclient "k8s.io/client-go/kubernetes"
)

// Reconciler reconciles a Potato object
type Reconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Config        *config.Config
	ClusterClient *kubernetesclient.Clientset
}

// +kubebuilder:rbac:groups=access.faros.sh,resources=request,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=access.faros.sh,resources=request/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=access.faros.sh,resources=request/finalizers,verbs=update

// Reconcile reconciles a Potato object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Include the clusterName from req.ObjectKey in the logger, similar to the namespace and name keys that are already
	// there.
	logger = logger.WithValues("clusterName", req.ClusterName)

	// Add the logical cluster to the context
	ctx = logicalcluster.WithCluster(ctx, logicalcluster.New(req.ClusterName))

	logger.Info("Getting Request")
	var original accessv1alpha1.Request
	if err := r.Get(ctx, req.NamespacedName, &original); err != nil {
		if errors.IsNotFound(err) {
			// Normal - was deleted
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	current := original.DeepCopy()

	original.Status.Message = "Success"

	logger.Info("KubeConfig status updated")

	if err := r.Status().Patch(ctx, &original, client.MergeFrom(current)); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&accessv1alpha1.Request{}).
		Complete(r)
}
