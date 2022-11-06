package workspaces

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	kcptenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

func (r *Reconciler) delete(ctx context.Context, logger logr.Logger, workspace *tenancyv1alpha1.Workspace) (ctrl.Result, error) {
	workspacePath := r.getWorkspaceName(workspace)
	cluster := logicalcluster.New(workspacePath)
	parent, exits := cluster.Parent()
	if !exits {
		return ctrl.Result{}, fmt.Errorf("parent cluster not found")
	}

	// global bindings
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getOrgClusterAccessName(workspace),
		},
	}

	result, err := r.deleteClusterRole(ctx, clusterRole, kcptenancyv1alpha1.RootCluster)
	if err != nil && !apierrors.IsNotFound(err) {
		return result, fmt.Errorf("failed to delete ClusterRole: %s", err)
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getOrgClusterAccessName(workspace),
		},
	}

	result, err = r.deleteClusterRoleBinding(ctx, clusterRoleBinding, kcptenancyv1alpha1.RootCluster)
	if err != nil && !apierrors.IsNotFound(err) {
		return result, fmt.Errorf("failed to delete ClusterRole: %s", err)
	}

	// localized bindings
	clusterRole = &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getRBACClusterAdminName(workspace),
		},
	}

	result, err = r.deleteClusterRole(ctx, clusterRole, parent)
	if err != nil && !apierrors.IsNotFound(err) {
		return result, fmt.Errorf("failed to delete ClusterRole: %s", err)
	}

	clusterRoleBinding = &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getRBACClusterAdminName(workspace),
		},
	}

	result, err = r.deleteClusterRoleBinding(ctx, clusterRoleBinding, parent)
	if err != nil && !apierrors.IsNotFound(err) {
		return result, fmt.Errorf("failed to delete ClusterRole: %s", err)
	}

	controllerutil.RemoveFinalizer(workspace, finalizerName)
	if err := r.Update(ctx, workspace); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) deleteClusterRole(ctx context.Context, clusterRole *rbacv1.ClusterRole, cluster logicalcluster.Name) (ctrl.Result, error) {
	return ctrl.Result{}, r.CoreClients.Cluster(cluster).RbacV1().ClusterRoles().Delete(ctx, clusterRole.Name, metav1.DeleteOptions{})
}

func (r *Reconciler) deleteClusterRoleBinding(ctx context.Context, clusterRoleBinding *rbacv1.ClusterRoleBinding, cluster logicalcluster.Name) (ctrl.Result, error) {
	return ctrl.Result{}, r.CoreClients.Cluster(cluster).RbacV1().ClusterRoleBindings().Delete(ctx, clusterRoleBinding.Name, metav1.DeleteOptions{})
}
