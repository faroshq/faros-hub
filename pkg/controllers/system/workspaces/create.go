package workspaces

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	kcptenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	kcptenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/bootstrap"
)

func (r *Reconciler) createOrUpdate(ctx context.Context, logger logr.Logger, workspace *tenancyv1alpha1.Workspace) (ctrl.Result, error) {
	// TODO: move to webhook
	if !controllerutil.ContainsFinalizer(workspace, finalizerName) {
		controllerutil.AddFinalizer(workspace, finalizerName)
		if err := r.Update(ctx, workspace); err != nil {
			return ctrl.Result{}, err
		}
		// requeue to ensure the finalizer is set before creating the resources
		return ctrl.Result{Requeue: true}, nil
	}

	workspaceOwnersReferences := []metav1.OwnerReference{{
		APIVersion: tenancyv1alpha1.SchemeGroupVersion.String(),
		Kind:       tenancyv1alpha1.WorkspaceKind,
		Name:       workspace.Name,
		UID:        workspace.UID,
	}}

	workspacePath := r.getWorkspaceName(workspace)
	cluster := logicalcluster.New(workspacePath)

	bootstraper, err := bootstrap.New(r.Config)
	if err != nil {
		return ctrl.Result{}, err
	}

	parent, exits := cluster.Parent()
	if !exits {
		return ctrl.Result{}, fmt.Errorf("parent cluster not found")
	}

	err = bootstraper.CreateWorkspace(ctx, parent.String())
	if err != nil {
		return ctrl.Result{}, err
	}
	rest, err := r.ClientFactory.GetChildWorkspaceRestConfig(ctx, workspacePath)
	if err != nil {
		return ctrl.Result{}, err
	}

	kcpClient, err := kcpclient.NewForConfig(rest)
	if err != nil {
		return ctrl.Result{}, err
	}

	ws := &kcptenancyv1beta1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:            workspace.Name,
			OwnerReferences: workspaceOwnersReferences,
		},
		Spec: kcptenancyv1beta1.WorkspaceSpec{
			Type: kcptenancyv1alpha1.ClusterWorkspaceTypeReference{
				Name: "faros",
				Path: "root",
			},
		},
	}

	_, err = kcpClient.TenancyV1beta1().Workspaces().Get(ctx, ws.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		logger.Error(err, "creating workspace", "workspace-name", workspace.Name)
		_, err := kcpClient.TenancyV1beta1().Workspaces().Create(ctx, ws, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create Workspace: %s", err)
		}
	case err == nil:
		// workspaces are not updatable
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to get the Workspace %s", err)
	}

	return ctrl.Result{}, nil
}
