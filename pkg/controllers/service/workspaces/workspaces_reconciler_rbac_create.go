package workspaces

import (
	"context"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	kcptenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/kcp-dev/logicalcluster/v3"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type workspaceRBACReconciler struct {
	getWorkspaceName                 func(w *tenancyv1alpha1.Workspace) string
	getOrgClusterAccessName          func(workspace *tenancyv1alpha1.Workspace) string
	getUserWithPrefixName            func(user string) string
	getRBACClusterAdminName          func(workspace *tenancyv1alpha1.Workspace) string
	getKCPWorkspace                  func(ctx context.Context, cluster logicalcluster.Path, name string) (*kcptenancyv1beta1.Workspace, error)
	createOrUpdateClusterRole        func(ctx context.Context, cluster logicalcluster.Path, clusterRole *rbacv1.ClusterRole) error
	createOrUpdateClusterRoleBinding func(ctx context.Context, cluster logicalcluster.Path, clusterRoleBinding *rbacv1.ClusterRoleBinding) error
}

func (r *workspaceRBACReconciler) reconcile(ctx context.Context, workspace *tenancyv1alpha1.Workspace) (reconcileStatus, error) {
	workspacePath := r.getWorkspaceName(workspace)
	clusterPath := logicalcluster.NewPath(workspacePath)
	rootPath := logicalcluster.NewPath("root")

	kcpWorkspace, err := r.getKCPWorkspace(ctx, clusterPath, workspace.Name)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// create global cluster role in root cluster to enable rbac
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getOrgClusterAccessName(workspace),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"tenancy.kcp.dev"},
				Resources:     []string{"workspaces/content"},
				Verbs:         []string{"access"},
				ResourceNames: []string{"faros-tenants"},
			},
		},
	}

	err = r.createOrUpdateClusterRole(ctx, rootPath, clusterRole)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// Role binding to enable the cluster role
	subjects := []rbacv1.Subject{}
	for _, member := range workspace.Spec.Members {
		subjects = append(subjects, rbacv1.Subject{
			Kind: rbacv1.UserKind,
			Name: r.getUserWithPrefixName(member),
		})
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.getOrgClusterAccessName(workspace),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     clusterRole.Name,
		},
		Subjects: subjects,
	}

	err = r.createOrUpdateClusterRoleBinding(ctx, rootPath, clusterRoleBinding)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// Create workspace dedicated roles for admin
	// TODO: Fow now 1 cluster role per workspace.
	// Optimize with merged of the rule
	clusterRole = &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:            r.getRBACClusterAdminName(workspace),
			OwnerReferences: getKCPWorkspaceOwnersReference(kcpWorkspace),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"tenancy.kcp.dev"},
				Resources:     []string{"workspaces/content"},
				Verbs:         []string{"admin"},
				ResourceNames: []string{workspace.Name},
			},
		},
	}

	err = r.createOrUpdateClusterRole(ctx, clusterPath, clusterRole)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	// Add binding for all workspace members requested
	subjects = []rbacv1.Subject{}
	for _, member := range workspace.Spec.Members {
		subjects = append(subjects, rbacv1.Subject{
			Kind: rbacv1.UserKind,
			Name: r.getUserWithPrefixName(member),
		})
	}

	clusterRoleBinding = &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:            r.getRBACClusterAdminName(workspace),
			OwnerReferences: getKCPWorkspaceOwnersReference(kcpWorkspace),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     clusterRole.Name,
		},
		Subjects: subjects,
	}

	err = r.createOrUpdateClusterRoleBinding(ctx, clusterPath, clusterRoleBinding)
	if err != nil {
		return reconcileStatusStopAndRequeue, err
	}

	conditions.MarkTrue(workspace, conditionsv1alpha1.ReadyCondition)

	return reconcileStatusContinue, nil
}
