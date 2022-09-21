package bootstrap

import (
	"context"
	"fmt"
	"os"
	"time"

	utilkubernetes "github.com/mjudeikis/kcp-example/pkg/util/kubernetes"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetesclient "k8s.io/client-go/kubernetes"
)

var defaultNamespace = "default"

func (b *bootstrap) createServiceAccount(ctx context.Context, workspace, serviceAccountName string) error {
	fmt.Printf("Creating service account %s in workspace %s \n", serviceAccountName, workspace)
	_, rest, err := b.getWorkspaceClient(ctx, b.kcpClient, b.rest, workspace)
	if err != nil {
		return err
	}

	client, err := kubernetesclient.NewForConfig(rest)
	if err != nil {
		return err
	}

	client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultNamespace,
		},
	}, metav1.CreateOptions{})

	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
	}

	var isUpdate bool
	current, err := client.CoreV1().ServiceAccounts(defaultNamespace).Get(ctx, sa.GetName(), metav1.GetOptions{})
	if err == nil {
		isUpdate = true
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if isUpdate {
		fmt.Printf("Updating serviceAccount - %s \n", sa.GetName())
		sa.SetResourceVersion(current.GetResourceVersion())
		_, err = client.CoreV1().ServiceAccounts(defaultNamespace).Update(ctx, sa, metav1.UpdateOptions{})
	} else {
		fmt.Printf("Creating serviceAccount - %s \n", sa.GetName())
		_, err = client.CoreV1().ServiceAccounts(defaultNamespace).Create(ctx, sa, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	return nil
}

func (b *bootstrap) createServiceAccountRoleBinding(ctx context.Context, workspace, serviceAccountName string) error {
	fmt.Printf("Creating service account %s role & rolebinding in workspace %s \n", serviceAccountName, workspace)
	_, rest, err := b.getWorkspaceClient(ctx, b.kcpClient, b.rest, workspace)
	if err != nil {
		return err
	}

	client, err := kubernetesclient.NewForConfig(rest)
	if err != nil {
		return err
	}

	{
		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name: serviceAccountName,
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"*"},
					Resources: []string{"*"},
					Verbs:     []string{"*"},
				},
			},
		}

		var isUpdate bool
		current, err := client.RbacV1().Roles(defaultNamespace).Get(ctx, role.GetName(), metav1.GetOptions{})
		if err == nil {
			isUpdate = true
		}
		if err != nil && !apierrors.IsNotFound(err) {
			return err
		}

		if isUpdate {
			fmt.Printf("Updating role - %s \n", role.GetName())
			role.SetResourceVersion(current.GetResourceVersion())
			_, err = client.RbacV1().Roles(defaultNamespace).Update(ctx, role, metav1.UpdateOptions{})
		} else {
			fmt.Printf("Creating role - %s \n", role.GetName())
			_, err = client.RbacV1().Roles(defaultNamespace).Create(ctx, role, metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}

	}
	{

		// Create Rolebinding

		roleBinding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: serviceAccountName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      serviceAccountName,
					Namespace: defaultNamespace,
				},
			},
			RoleRef: rbacv1.RoleRef{
				Kind: "Role",
				Name: serviceAccountName,
			},
		}

		var isUpdate bool
		current, err := client.RbacV1().RoleBindings(defaultNamespace).Get(ctx, roleBinding.GetName(), metav1.GetOptions{})
		if err == nil {
			isUpdate = true
		}
		if err != nil && !apierrors.IsNotFound(err) {
			return err
		}

		if isUpdate {
			fmt.Printf("Updating roleBinding - %s \n", roleBinding.GetName())
			roleBinding.SetResourceVersion(current.GetResourceVersion())
			_, err = client.RbacV1().RoleBindings(defaultNamespace).Update(ctx, roleBinding, metav1.UpdateOptions{})
		} else {
			fmt.Printf("Creating roleBinding - %s \n", roleBinding.GetName())
			_, err = client.RbacV1().RoleBindings(defaultNamespace).Create(ctx, roleBinding, metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *bootstrap) createServiceAccountKubeconfig(ctx context.Context, workspace, serviceAccountName, path string) error {
	fmt.Printf("Creating service account kubeconfig %s for workspace %s \n", serviceAccountName, workspace)
	_, rest, err := b.getWorkspaceClient(ctx, b.kcpClient, b.rest, workspace)
	if err != nil {
		return err
	}

	client, err := kubernetesclient.NewForConfig(rest)
	if err != nil {
		return err
	}

	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
	}

	var found bool
	var current *corev1.ServiceAccount
	for !found {
		current, err = client.CoreV1().ServiceAccounts(defaultNamespace).Get(ctx, sa.GetName(), metav1.GetOptions{})
		if err == nil && len(current.Secrets) > 0 {
			found = true
		}
		time.Sleep(1 * time.Second)
	}

	secretName := current.Secrets[0].Name
	secret, err := client.CoreV1().Secrets(defaultNamespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	token := secret.Data["token"]
	namespace := secret.Namespace

	server := rest.Host

	kc, err := utilkubernetes.MakeKubeconfig(server, namespace, string(token))
	if err != nil {
		return err
	}

	err = os.WriteFile(path, kc, 0644)
	if err != nil {
		return err
	}

	return nil
}
