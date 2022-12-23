package workspaces

import (
	"context"
	"fmt"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	kcptenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	"github.com/kcp-dev/logicalcluster/v3"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type reconcileStatus int

const (
	reconcileStatusStopAndRequeue reconcileStatus = iota
	reconcileStatusContinue
	reconcileStatusError
)

type reconciler interface {
	reconcile(ctx context.Context, workspace *tenancyv1alpha1.Workspace) (reconcileStatus, error)
}

func (c *Controller) reconcile(ctx context.Context, ws *tenancyv1alpha1.Workspace) (bool, error) {
	var reconcilers []reconciler
	createReconcilers := []reconciler{
		&finalizerAddReconciler{ // must be first
			getFinalizerName: func() string {
				return finalizerName
			},
		},
		&kcpWorkspaceReconciler{ // must be second
			createKCPWorkspace: c.bootstraper.CreateWorkspace,
			getWorkspaceName: func(w *tenancyv1alpha1.Workspace) string {
				return getWorkspaceName(c.config, w)
			},
			createFarosWorkspace: func(ctx context.Context, cluster logicalcluster.Path, workspace *tenancyv1alpha1.Workspace) error {
				return c.createFarosWorkspace(ctx, cluster, workspace)
			},
		},
		&workspaceRBACReconciler{ // must be after kcpWorkspaceReconciler
			getWorkspaceName: func(w *tenancyv1alpha1.Workspace) string {
				return getWorkspaceName(c.config, w)
			},
			getKCPWorkspace: func(ctx context.Context, cluster logicalcluster.Path, name string) (*kcptenancyv1beta1.Workspace, error) {
				return c.kcpClientSet.Cluster(cluster).TenancyV1beta1().Workspaces().Get(ctx, name, metav1.GetOptions{})
			},
			getOrgClusterAccessName: func(w *tenancyv1alpha1.Workspace) string {
				return getOrgClusterAccessName(c.config, w)
			},
			createOrUpdateClusterRole: func(ctx context.Context, cluster logicalcluster.Path, role *rbacv1.ClusterRole) error {
				return createOrUpdateClusterRole(ctx, c.coreClientSet.Cluster(cluster), role)
			},
			createOrUpdateClusterRoleBinding: func(ctx context.Context, cluster logicalcluster.Path, roleBinding *rbacv1.ClusterRoleBinding) error {
				return createOrUpdateClusterRoleBinding(ctx, c.coreClientSet.Cluster(cluster), roleBinding)
			},
			getUserWithPrefixName: func(name string) string {
				return getUserWithPrefixName(c.config, name)
			},
			getRBACClusterAdminName: func(w *tenancyv1alpha1.Workspace) string {
				return getRBACClusterAdminName(w)
			},
		},
	}

	deleteReconcilers := []reconciler{
		&workspaceRBACDeleteReconciler{
			getWorkspaceName: func(w *tenancyv1alpha1.Workspace) string {
				return getWorkspaceName(c.config, w)
			},
			getOrgClusterAccessName: func(w *tenancyv1alpha1.Workspace) string {
				return getOrgClusterAccessName(c.config, w)
			},
			deleteClusterRole: func(ctx context.Context, cluster logicalcluster.Path, role *rbacv1.ClusterRole) error {
				return c.coreClientSet.Cluster(cluster).RbacV1().ClusterRoles().Delete(ctx, role.Name, metav1.DeleteOptions{})
			},
			deleteClusterRoleBinding: func(ctx context.Context, cluster logicalcluster.Path, roleBinding *rbacv1.ClusterRoleBinding) error {
				return c.coreClientSet.Cluster(cluster).RbacV1().ClusterRoleBindings().Delete(ctx, roleBinding.Name, metav1.DeleteOptions{})
			},
			getUserWithPrefixName: func(name string) string {
				return getUserWithPrefixName(c.config, name)
			},
			getRBACClusterAdminName: func(w *tenancyv1alpha1.Workspace) string {
				return getRBACClusterAdminName(w)
			},
		},
		&kcpWorkspaceDeleteReconciler{
			deleteFarosWorkspace: func(ctx context.Context, cluster logicalcluster.Path, workspace *tenancyv1alpha1.Workspace) error {
				return c.deleteFarosWorkspace(ctx, cluster, workspace)
			},
			getWorkspaceName: func(w *tenancyv1alpha1.Workspace) string {
				return getWorkspaceName(c.config, w)
			},
		},
		&finalizerRemoveReconciler{
			getFinalizerName: func() string {
				return finalizerName
			},
		},
	}

	if !ws.DeletionTimestamp.IsZero() { //delete
		reconcilers = deleteReconcilers
	} else { //create or update
		reconcilers = createReconcilers
	}

	var errs []error

	requeue := false
	for _, r := range reconcilers {
		var err error
		var status reconcileStatus
		status, err = r.reconcile(ctx, ws)
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

func (c *Controller) createFarosWorkspace(ctx context.Context, cluster logicalcluster.Path, workspace *tenancyv1alpha1.Workspace) error {
	logger := klog.FromContext(ctx)

	ws := &kcptenancyv1beta1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name: workspace.Name,
		},
		Spec: kcptenancyv1beta1.WorkspaceSpec{
			Type: kcptenancyv1beta1.WorkspaceTypeReference{
				Name: "faros",
				Path: "root",
			},
		},
	}

	kcpWorkspace, err := c.kcpClientSet.Cluster(cluster).TenancyV1beta1().Workspaces().Get(ctx, ws.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		logger.Info("creating workspace", "workspace-name", workspace.Name)
		kcpWorkspace, err = c.kcpClientSet.Cluster(cluster).TenancyV1beta1().Workspaces().Create(ctx, ws, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create Workspace: %s", err)
		}
	case err == nil:
		// workspaces are not updatable, but we need to deal with all the stuff bellow
	default:
		return fmt.Errorf("failed to get the Workspace %s", err)
	}

	workspace.Status.WorkspaceURL = kcpWorkspace.Status.URL
	return nil
}

func (c *Controller) deleteFarosWorkspace(ctx context.Context, cluster logicalcluster.Path, workspace *tenancyv1alpha1.Workspace) error {
	ws := &kcptenancyv1beta1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name: workspace.Name,
		},
	}

	return c.kcpClientSet.Cluster(cluster).TenancyV1beta1().Workspaces().Delete(ctx, ws.Name, metav1.DeleteOptions{})
}

func createOrUpdateClusterRole(ctx context.Context, coreClients kubernetes.Interface, clusterRole *rbacv1.ClusterRole) error {
	originalOwners := clusterRole.OwnerReferences

	currentClusterRole, err := coreClients.RbacV1().ClusterRoles().Get(ctx, clusterRole.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err := coreClients.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create the ClusterRole %s", err)
		}
	case err == nil:
		currentClusterRole.Rules = clusterRole.Rules
		currentClusterRole.ResourceVersion = ""
		currentClusterRole.OwnerReferences = mergeOwnerReference(currentClusterRole.OwnerReferences, originalOwners)
		_, err := coreClients.RbacV1().ClusterRoles().Update(ctx, currentClusterRole, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update the ClusterRole %s", err)
		}
	default:
		return fmt.Errorf("failed to create the ClusterRole %s", err)
	}

	return nil
}

func createOrUpdateClusterRoleBinding(ctx context.Context, coreClients kubernetes.Interface, clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	originalOwners := clusterRoleBinding.OwnerReferences

	currentClusterRoleBinding, err := coreClients.RbacV1().ClusterRoleBindings().Get(ctx, clusterRoleBinding.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err := coreClients.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create the ClusterRoleBindings %s", err)
		}
	case err == nil:
		currentClusterRoleBinding.RoleRef = clusterRoleBinding.RoleRef
		currentClusterRoleBinding.Subjects = clusterRoleBinding.Subjects
		currentClusterRoleBinding.ResourceVersion = ""
		currentClusterRoleBinding.OwnerReferences = mergeOwnerReference(currentClusterRoleBinding.OwnerReferences, originalOwners)
		_, err := coreClients.RbacV1().ClusterRoleBindings().Update(ctx, currentClusterRoleBinding, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update the ClusterRoleBindings %s", err)
		}
	default:
		return fmt.Errorf("failed to create the ClusterRoleBindings %s", err)
	}

	return nil
}
