package access

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/kcp-dev/logicalcluster/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	accessv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/access/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/util/kubeconfig"
	utilkubernetes "github.com/faroshq/faros-hub/pkg/util/kubernetes"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
)

var kubeRootCA = "kube-root-ca.crt"

// Reconciler reconciles a Potato object
type Reconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Config        *config.Config
	ClientFactory utilkubernetes.ClientFactory
	CoreClients   kubernetes.ClusterInterface
}

// +kubebuilder:rbac:groups=access.faros.sh,resources=request,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=access.faros.sh,resources=request/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=access.faros.sh,resources=request/finalizers,verbs=update

// Reconcile reconciles a Potato object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Include the clusterName from req.ObjectKey in the logger, similar to the namespace and name keys that are already
	// there.
	logger = logger.WithValues("clusterName", req.ClusterName).WithValues("namespace", req.Namespace).WithValues("name", req.Name)

	// Add the logical cluster to the context
	ctx = logicalcluster.WithCluster(ctx, logicalcluster.New(req.ClusterName))

	logger.Info("Getting Request")
	var request accessv1alpha1.Request
	if err := r.Get(ctx, req.NamespacedName, &request); err != nil {
		if errors.IsNotFound(err) {
			// Normal - was deleted
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	kcpClient, err := r.ClientFactory.GetRootKCPClient()
	if err != nil {
		return ctrl.Result{}, err
	}

	requestCopy := request.DeepCopy()

	// check if sync target exists and is ready
	synctarget, err := kcpClient.Cluster(logicalcluster.New(req.ClusterName)).WorkloadV1alpha1().SyncTargets().Get(ctx, request.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		conditions.MarkFalse(requestCopy, conditionsv1alpha1.ReadyCondition, "SyncTargetNotFound", conditionsv1alpha1.ConditionSeverityError, err.Error())
		if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, err
		}
		return ctrl.Result{RequeueAfter: time.Second * 30}, err
	}

	if !conditions.IsTrue(synctarget, conditionsv1alpha1.ReadyCondition) {
		conditions.MarkFalse(requestCopy, conditionsv1alpha1.ReadyCondition, "SyncTargetNotReady", conditionsv1alpha1.ConditionSeverityError, "SyncTarget is not ready")
		if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, err
		}
		return ctrl.Result{RequeueAfter: time.Second * 30}, err
	}

	requestSecret := corev1.Secret{}
	requestSecret.Name = request.Name
	requestSecret.Namespace = request.Namespace
	err = r.Get(ctx, client.ObjectKey{Name: requestSecret.Name, Namespace: requestSecret.Namespace}, &requestSecret)
	if err != nil {
		if !errors.IsNotFound(err) {
			conditions.MarkFalse(requestCopy, conditionsv1alpha1.ReadyCondition, "FailedToGetRequestSecret", conditionsv1alpha1.ConditionSeverityError, err.Error())
			if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
	}
	exists := err == nil
	if !exists {
		logger.Info("Creating request temporary credentials secrets secret")
		var caConfigMap corev1.ConfigMap

		err := r.Client.Get(ctx, client.ObjectKey{Name: kubeRootCA, Namespace: request.Namespace}, &caConfigMap)
		if err != nil {
			conditions.MarkFalse(requestCopy, conditionsv1alpha1.ReadyCondition, "FailedToGetRootCASecret", conditionsv1alpha1.ConditionSeverityError, err.Error())
			if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}

		caCrt, ok := caConfigMap.Data["ca.crt"]
		if !ok {
			conditions.MarkFalse(requestCopy, conditionsv1alpha1.ReadyCondition, "FailedToGetRootCASecret", conditionsv1alpha1.ConditionSeverityError, "ca.crt not found in configmap")
			if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}

		token := uuid.New().String()
		kubeconfig, err := r.generateKubeConfig(ctx, requestCopy, req.ClusterName, token, caCrt)
		if err != nil {
			conditions.MarkFalse(requestCopy, conditionsv1alpha1.ReadyCondition, "FailedToGenerateKubeconfig", conditionsv1alpha1.ConditionSeverityError, err.Error())
			if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}

		requestSecret.Data = map[string][]byte{
			"token":      []byte(token),
			"kubeconfig": kubeconfig,
		}

		if err := r.Create(ctx, &requestSecret); err != nil {
			conditions.MarkFalse(requestCopy, conditionsv1alpha1.ReadyCondition, "FailedToCreateRequestSecret", conditionsv1alpha1.ConditionSeverityError, err.Error())
			if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
	}

	conditions.MarkTrue(requestCopy, conditionsv1alpha1.ReadyCondition)
	requestCopy.Status.SyncTarget = synctarget.Name
	logger.Info("KubeConfig status updated")
	if err := r.Status().Patch(ctx, requestCopy, client.MergeFrom(&request)); err != nil {
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

// https://host/services/faros-tunnels/clusters/<ws>/apis/access.faros.sh/v1alpha1/clusters/<name>/connect establish reverse connections and queue them so it can be consumed by the dialer
// https://host/services/faros-tunnels/clusters/<ws>/apis/access.faros.sh/v1alpha1/clusters/<name>/proxy/{path} proxies the {path} through the reverse connection identified by the cluster and syncer name
func (r *Reconciler) generateKubeConfig(ctx context.Context, request *accessv1alpha1.Request, cluster, token, cacrt string) ([]byte, error) {
	path := fmt.Sprintf("/services/faros-tunnels/clusters/%s/apis/access.faros.sh/v1alpha1/namespaces/%s/access/%s/proxy",
		cluster, request.Namespace, request.Name)

	u, err := url.Parse(r.Config.RootRestConfig.Host)
	if err != nil {
		return nil, err
	}
	u.Path = ""

	server := u.String() + path
	// TODO: inject CA so we can verify the server
	// Currently we are using insecure-skip-tls-verify

	return kubeconfig.MakeKubeconfig(server, token, cacrt)
}
