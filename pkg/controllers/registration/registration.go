package registration

import (
	"context"
	"fmt"
	"time"

	"github.com/kcp-dev/logicalcluster/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
	workloadv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/workload/v1alpha1"
)

// Reconciler reconciles a Potato object
type Reconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Config        *config.Config
	ClientFactory utilkubernetes.ClientFactory
	CoreClients   kubernetes.ClusterInterface
}

// +kubebuilder:rbac:groups=edge.faros.sh,resources=registration,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=edge.faros.sh,resources=registration/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=edge.faros.sh,resources=registration/finalizers,verbs=update

// Reconcile reconciles a Potato object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Include the clusterName from req.ObjectKey in the logger, similar to the namespace and name keys that are already
	// there.
	logger = logger.WithValues("clusterName", req.ClusterName).WithValues("namespace", req.Namespace).WithValues("name", req.Name)

	// Add the logical cluster to the context
	ctx = logicalcluster.WithCluster(ctx, logicalcluster.New(req.ClusterName))

	logger.Info("Getting registration")
	var registration edgev1alpha1.Registration
	if err := r.Get(ctx, req.NamespacedName, &registration); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	registrationOwnersReferences := []metav1.OwnerReference{{
		APIVersion: workloadv1alpha1.SchemeGroupVersion.String(),
		Kind:       edgev1alpha1.RegistrationKind,
		Name:       registration.Name,
		UID:        registration.UID,
	}}

	var sa *corev1.ServiceAccount
	err := r.Get(ctx, req.NamespacedName, sa)
	switch {
	case apierrors.IsNotFound(err):
		logger.Error(err, "Creating service account", req.Name)
		err := r.Create(ctx, &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:            req.Name,
				Namespace:       req.Namespace,
				OwnerReferences: registrationOwnersReferences,
			},
		}, &client.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create ServiceAccount: %s", err)
		}
	case err == nil:
		// service account already exist, merge owner references
		sa.OwnerReferences = mergeOwnerReference(sa.OwnerReferences, registrationOwnersReferences)

		logger.Info("Updating service account %s", sa.Name)
		err = r.Patch(ctx, sa, client.MergeFrom(sa.DeepCopy()), &client.PatchOptions{})
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to patch ServiceAccount %s", err)
		}
	default:
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to get the ServiceAccount %s", err)
	}

	return ctrl.Result{}, nil
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
