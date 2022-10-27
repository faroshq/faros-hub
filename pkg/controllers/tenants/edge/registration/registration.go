package registration

import (
	"context"

	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/kcp-dev/logicalcluster/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
)

var finalizerName = "registration.edge.faros.sh/finalizer"

// Reconciler reconciles a Registration object
type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=edge.faros.sh,resources=registrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=edge.faros.sh,resources=registrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=edge.faros.sh,resources=registrations/finalizers,verbs=update

// Reconcile reconciles a Registration object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Include the clusterName from req.ObjectKey in the logger, similar to the namespace and name keys that are already
	// there.
	logger = logger.WithValues("clusterName", req.ClusterName).WithValues("namespace", req.Namespace).WithValues("name", req.Name)

	// Add the logical cluster to the context
	ctx = logicalcluster.WithCluster(ctx, logicalcluster.New(req.ClusterName))

	var registration edgev1alpha1.Registration
	if err := r.Get(ctx, req.NamespacedName, &registration); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var result ctrl.Result
	var err error
	if registration.DeletionTimestamp.IsZero() {
		result, err = r.createOrUpdate(ctx, logger, registration.DeepCopy())
	} else {
		result, err = r.delete(ctx, logger, registration.DeepCopy())
	}
	if err != nil {
		registrationCopy := registration.DeepCopy()
		conditions.MarkFalse(
			registrationCopy,
			conditionsv1alpha1.ReadyCondition,
			err.Error(),
			conditionsv1alpha1.ConditionSeverityError,
			"Error configuring Registration: %v",
			err,
		)
		if err := r.Status().Patch(ctx, registrationCopy, client.MergeFrom(&registration)); err != nil {
			return result, err
		}
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&edgev1alpha1.Registration{}).
		Complete(r)
}

// mergeOwnerReference: merge a slice of ownerReference with a given ownerReferences
func mergeOwnerReference(ownerReferences, newOwnerReferences []metav1.OwnerReference) []metav1.OwnerReference {
	var merged []metav1.OwnerReference

	merged = append(merged, ownerReferences...)

	for _, ownerReference := range newOwnerReferences {
		found := false
		for _, mergedOwnerReference := range merged {
			if mergedOwnerReference.UID == ownerReference.UID {
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, ownerReference)
		}
	}

	return merged

}
