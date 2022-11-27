package workspaces

import (
	"context"
	"fmt"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	kcptenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v2"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type workspaceRBACDeleteReconciler struct {
	getWorkspaceName         func(w *tenancyv1alpha1.Workspace) string
	getOrgClusterAccessName  func(workspace *tenancyv1alpha1.Workspace) string
	getUserWithPrefixName    func(user string) string
	getRBACClusterAdminName  func(workspace *tenancyv1alpha1.Workspace) string
	deleteClusterRole        func(ctx context.Context, cluster logicalcluster.Name, clusterRole *rbacv1.ClusterRole) error
	deleteClusterRoleBinding func(ctx context.Context, cluster logicalcluster.Name, clusterRoleBinding *rbacv1.ClusterRoleBinding) error
}

func (r *workspaceRBACDeleteReconciler) reconcile(ctx context.Context, workspace *tenancyv1alpha1.Workspace) (reconcileStatus, error) {
	workspacePath := r.getWorkspaceName(workspace)
	cluster := logicalcluster.New(workspacePath)

	parent, exits := cluster.Parent()
	if !exits {
		return reconcileStatusError, fmt.Errorf("parent cluster not found")
	}

	// global bindings
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getOrgClusterAccessName(workspace),
		},
	}

	err := r.deleteClusterRole(ctx, kcptenancyv1alpha1.RootCluster, clusterRole)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getOrgClusterAccessName(workspace),
		},
	}

	err = r.deleteClusterRoleBinding(ctx, kcptenancyv1alpha1.RootCluster, clusterRoleBinding)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// localized bindings
	clusterRole = &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getRBACClusterAdminName(workspace),
		},
	}

	err = r.deleteClusterRole(ctx, parent, clusterRole)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	clusterRoleBinding = &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getRBACClusterAdminName(workspace),
		},
	}

	err = r.deleteClusterRoleBinding(ctx, parent, clusterRoleBinding)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	return reconcileStatusContinue, nil
}
