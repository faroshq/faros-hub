package pluginrequests

import (
	"context"
	"fmt"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/models"
	v1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
	kcpclientset "github.com/kcp-dev/kcp/pkg/client/clientset/versioned/cluster"
	"github.com/kcp-dev/logicalcluster/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type reconcileStatus int

const (
	reconcileStatusStopAndRequeue reconcileStatus = iota
	reconcileStatusContinue
	reconcileStatusError
)

type reconciler interface {
	reconcile(ctx context.Context, cluster logicalcluster.Name, request *pluginsv1alpha1.Request) (reconcileStatus, error)
}

func (c *Controller) reconcile(ctx context.Context, cluster logicalcluster.Name, request *pluginsv1alpha1.Request) (bool, error) {
	var reconcilers []reconciler
	createReconcilers := []reconciler{
		&finalizerAddReconciler{ // must be first
			getFinalizerName: func() string {
				return finalizerName
			},
		},
		&requestCreateReconciler{
			getPlugins: func() models.PluginsList {
				return c.plugins
			},
			createAPIExport: func(ctx context.Context, destinationCluster logicalcluster.Name, pluginVersion, pluginName string) error {
				sourceCluster := logicalcluster.New(c.config.ControllersPluginsWorkspace)
				sourceName := fmt.Sprintf("%s.%s", pluginVersion, pluginName)
				destinationName := pluginName

				return createAPIBinding(ctx, c.kcpClientSet, sourceCluster, destinationCluster, sourceName, destinationName)
			},
		},
	}

	deleteReconcilers := []reconciler{
		&finalizerRemoveReconciler{
			getFinalizerName: func() string {
				return finalizerName
			},
		},
		&requestDeleteReconciler{},
	}

	if !request.DeletionTimestamp.IsZero() { //delete
		reconcilers = deleteReconcilers
	} else { //create or update
		reconcilers = createReconcilers
	}

	var errs []error

	requeue := false
	for _, r := range reconcilers {
		var err error
		var status reconcileStatus
		status, err = r.reconcile(ctx, cluster, request)
		if err != nil {
			errs = append(errs, err)
		}
		if status == reconcileStatusStopAndRequeue {
			requeue = true
			break
		}
	}

	return requeue, utilerrors.NewAggregate(errs)
}

func createAPIBinding(ctx context.Context, kcpClient kcpclientset.ClusterInterface, sourceCluster, destinationCluster logicalcluster.Name, sourceName, destinationName string) error {
	apiExport, err := kcpClient.Cluster(sourceCluster).ApisV1alpha1().APIExports().Get(ctx, sourceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	apiBinding := &v1alpha1.APIBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: destinationName,
		},
		Spec: v1alpha1.APIBindingSpec{
			Reference: v1alpha1.ExportReference{
				Workspace: &v1alpha1.WorkspaceExportReference{
					Path:       sourceCluster.String(),
					ExportName: apiExport.Name,
				},
			},
		},
	}

	current, err := kcpClient.Cluster(destinationCluster).ApisV1alpha1().APIBindings().Get(ctx, apiBinding.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err = kcpClient.Cluster(destinationCluster).ApisV1alpha1().APIBindings().Create(ctx, apiBinding, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create the APIBinding %s", err)
		}
	case err == nil:
		current.Spec = apiBinding.Spec
		//current.ResourceVersion = ""
		_, err = kcpClient.Cluster(destinationCluster).ApisV1alpha1().APIBindings().Update(ctx, current, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update the APIBinding %s", err)
		}
	default:
		return fmt.Errorf("failed to create the APIBinding %s", err)
	}
	return nil
}
