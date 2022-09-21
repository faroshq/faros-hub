/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/mjudeikis/kcp-example/pkg/config"
	potatoesv1alpha1 "github.com/mjudeikis/kcp-example/pkg/controllers/potatoes/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	kubernetesclient "k8s.io/client-go/kubernetes"
)

// Reconciler reconciles a Potato object
type Reconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Config        *config.Config
	ClusterClient *kubernetesclient.Clientset

	lock sync.Mutex
}

// +kubebuilder:rbac:groups=faros.sh,resources=potatoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=faros.sh,resources=potatoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=faros.sh,resources=potatoes/finalizers,verbs=update

// Reconcile reconciles a Potato object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Include the clusterName from req.ObjectKey in the logger, similar to the namespace and name keys that are already
	// there.
	logger = logger.WithValues("clusterName", req.ClusterName)

	// Add the logical cluster to the context
	ctx = logicalcluster.WithCluster(ctx, logicalcluster.New(req.ClusterName))

	logger.Info("Getting Potato request")
	var original potatoesv1alpha1.Potato
	if err := r.Get(ctx, req.NamespacedName, &original); err != nil {
		if errors.IsNotFound(err) {
			// Normal - was deleted
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	current := original.DeepCopy()

	logger.Info("Getting Potatos from warehouse")
	err := r.GetPotatoes(original.Spec.Request)
	if err != nil {
		logger.Error(err, "Failed to get potatoes")
		original.Status.Message = err.Error()
	} else {
		original.Status.Message = "Success"
		original.Status.Total = original.Spec.Request
	}

	logger.Info("Patching potato status to store total potato count in the remote cluster")

	if err := r.Status().Patch(ctx, &original, client.MergeFrom(current)); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

var (
	configMapName = "warehouse"
	keyName       = "potatoes"
)

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: os.Getenv("POD_NAMESPACE"),
		},
		Data: map[string]string{
			keyName: strconv.FormatInt(r.Config.Server.ControllerPotatoesCount, 10),
		},
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	_, err := r.ClusterClient.CoreV1().ConfigMaps(cm.GetNamespace()).Create(ctx, &cm, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	if errors.IsAlreadyExists(err) {
		_current, err := r.ClusterClient.CoreV1().ConfigMaps(cm.GetNamespace()).Get(ctx, configMapName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		current := _current.DeepCopy()
		current.Data[keyName] = strconv.FormatInt(r.Config.Server.ControllerPotatoesCount, 10)
		_, err = r.ClusterClient.CoreV1().ConfigMaps(cm.GetNamespace()).Update(ctx, current, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&potatoesv1alpha1.Potato{}).
		Complete(r)
}

func (r *Reconciler) GetPotatoes(count int64) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: os.Getenv("POD_NAMESPACE"),
		},
	}
	_current, err := r.ClusterClient.CoreV1().ConfigMaps(cm.GetNamespace()).Get(ctx, cm.GetName(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	current := _current.DeepCopy()

	bananas, err := strconv.ParseInt(current.Data[keyName], 10, 64)
	if err != nil {
		return err
	}
	if bananas < count {
		return errors.NewBadRequest("not enough potatoes in the warehouse")
	} else {
		bananas -= count
		current.Data[keyName] = strconv.FormatInt(bananas, 10)
		_, err = r.ClusterClient.CoreV1().ConfigMaps(cm.GetNamespace()).Update(ctx, current, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
