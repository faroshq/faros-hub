package workspaces

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	kcptenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	kcptenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v2"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		APIVersion:         tenancyv1alpha1.SchemeGroupVersion.String(),
		Kind:               tenancyv1alpha1.WorkspaceKind,
		Name:               workspace.Name,
		BlockOwnerDeletion: pointer.BoolPtr(true),
		Controller:         pointer.BoolPtr(true),
		UID:                workspace.UID,
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
		// workspaces are not updatable, but we need to deal with all the stuff bellow
	default:
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to get the Workspace %s", err)
	}

	// create global cluster role in root cluster to enable rbac
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:            r.getOrgClusterAccessName(workspace),
			OwnerReferences: workspaceOwnersReferences,
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

	result, err := r.createOrUpdateClusterRole(ctx, clusterRole, kcptenancyv1alpha1.RootCluster, workspaceOwnersReferences)
	if err != nil {
		return result, err
	}

	// Role binding to enable the cluster role
	subjects := []rbacv1.Subject{}
	for _, member := range workspace.Spec.Members {
		subjects = append(subjects, rbacv1.Subject{
			Kind: rbacv1.UserKind,
			Name: fmt.Sprintf("%s:%s", r.Config.OIDCUserPrefix, member),
		})
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:            r.getOrgClusterAccessName(workspace),
			OwnerReferences: workspaceOwnersReferences,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     clusterRole.Name,
		},
		Subjects: subjects,
	}

	result, err = r.createOrUpdateClusterRoleBinding(ctx, clusterRoleBinding, kcptenancyv1alpha1.RootCluster, workspaceOwnersReferences)
	if err != nil {
		return result, err
	}

	// Create workspace dedicated roles for admin
	// TODO: Fow now 1 cluster role per workspace.
	// Optimize with merged of the rule
	clusterRole = &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:            r.getRBACClusterAdminName(workspace),
			OwnerReferences: workspaceOwnersReferences,
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

	result, err = r.createOrUpdateClusterRole(ctx, clusterRole, parent, workspaceOwnersReferences)
	if err != nil {
		return result, err
	}

	// Add binding for all workspace members requested
	subjects = []rbacv1.Subject{}
	for _, member := range workspace.Spec.Members {
		subjects = append(subjects, rbacv1.Subject{
			Kind: rbacv1.UserKind,
			Name: fmt.Sprintf("%s:%s", r.Config.OIDCUserPrefix, member),
		})
	}

	clusterRoleBinding = &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:            r.getRBACClusterAdminName(workspace),
			OwnerReferences: workspaceOwnersReferences,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     clusterRole.Name,
		},
		Subjects: subjects,
	}

	result, err = r.createOrUpdateClusterRoleBinding(ctx, clusterRoleBinding, parent, workspaceOwnersReferences)
	if err != nil {
		return result, err
	}

	patch := client.MergeFrom(workspace.DeepCopy())
	conditions.MarkTrue(workspace, conditionsv1alpha1.ReadyCondition)

	if err := r.Status().Patch(ctx, workspace, patch); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
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

func (r *Reconciler) createOrUpdateClusterRole(ctx context.Context, clusterRole *rbacv1.ClusterRole, cluster logicalcluster.Name, owners []metav1.OwnerReference) (ctrl.Result, error) {
	currentClusterRole, err := r.CoreClients.Cluster(cluster).RbacV1().ClusterRoles().Get(ctx, clusterRole.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err := r.CoreClients.Cluster(cluster).RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create the ClusterRole %s", err)
		}
	case err == nil:
		currentClusterRole.Rules = clusterRole.Rules
		currentClusterRole.ResourceVersion = ""
		currentClusterRole.OwnerReferences = mergeOwnerReference(clusterRole.OwnerReferences, owners)
		_, err := r.CoreClients.Cluster(cluster).RbacV1().ClusterRoles().Update(ctx, currentClusterRole, metav1.UpdateOptions{})
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to update the ClusterRole %s", err)
		}
	default:
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create the ClusterRole %s", err)
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) createOrUpdateClusterRoleBinding(ctx context.Context, clusterRoleBinding *rbacv1.ClusterRoleBinding, cluster logicalcluster.Name, owners []metav1.OwnerReference) (ctrl.Result, error) {
	currentClusterRoleBinding, err := r.CoreClients.Cluster(cluster).RbacV1().ClusterRoleBindings().Get(ctx, clusterRoleBinding.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		_, err := r.CoreClients.Cluster(cluster).RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create the ClusterRoleBindings %s", err)
		}
	case err == nil:
		currentClusterRoleBinding.RoleRef = clusterRoleBinding.RoleRef
		currentClusterRoleBinding.Subjects = clusterRoleBinding.Subjects
		currentClusterRoleBinding.ResourceVersion = ""
		currentClusterRoleBinding.OwnerReferences = mergeOwnerReference(clusterRoleBinding.OwnerReferences, owners)
		_, err := r.CoreClients.Cluster(cluster).RbacV1().ClusterRoleBindings().Update(ctx, clusterRoleBinding, metav1.UpdateOptions{})
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to update the ClusterRoleBindings %s", err)
		}
	default:
		return ctrl.Result{RequeueAfter: time.Second * 30}, fmt.Errorf("failed to create the ClusterRoleBindings %s", err)
	}

	return ctrl.Result{}, nil
}
