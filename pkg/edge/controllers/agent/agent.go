package agent

import (
	"context"

	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/config"
)

// Reconciler reconciles a Potato object
type Reconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Config      *config.AgentConfig
	FarosClient farosclient.Interface
}

// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=edge.faros.sh,resources=agent/finalizers,verbs=update

// Reconcile reconciles a Edge object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if req.Name != r.Config.Name && req.Namespace != r.Config.Namespace {
		return ctrl.Result{}, nil
	}

	// TODO: For some reason dynamic client from controller-runtime can't get if we scope it to a namespace
	agent, err := r.FarosClient.EdgeV1alpha1().Agents(r.Config.Namespace).Get(ctx, r.Config.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	agentCopy := agent.DeepCopy()
	conditions.MarkTrue(agentCopy, conditionsv1alpha1.ReadyCondition)

	_, err = r.FarosClient.EdgeV1alpha1().Agents(r.Config.Namespace).UpdateStatus(ctx, agentCopy, metav1.UpdateOptions{})
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&edgev1alpha1.Agent{}).
		Complete(r)
}
