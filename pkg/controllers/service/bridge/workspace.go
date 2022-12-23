package bridge

import (
	"context"

	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/faroshq/faros-hub/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (b *bridge) createWorkspace(ctx context.Context, workspaceID string) error {
	w, err := b.store.GetWorkspace(ctx, models.Workspace{
		ID: workspaceID,
	})
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}
	if err == store.ErrRecordNotFound { // gone, nothing to do
		return nil
	}

	workspace := w.Workspace
	if workspace.Labels == nil {
		workspace.Labels = map[string]string{}
	}
	workspace.Labels[models.LabelIdentifier] = workspaceID

	_, err = b.farosClientSet.TenancyV1alpha1().Workspaces().Create(ctx, &workspace, metav1.CreateOptions{})
	return err
}

func (b *bridge) deleteWorkspace(ctx context.Context, workspaceID string) error {
	err := b.farosClientSet.TenancyV1alpha1().Workspaces().DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: models.LabelIdentifier + "=" + workspaceID,
	})
	return err
}

func (b *bridge) updateWorkspace(ctx context.Context, workspaceID string) error {
	w, err := b.store.GetWorkspace(ctx, models.Workspace{
		ID: workspaceID,
	})
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}
	if err == store.ErrRecordNotFound { // gone, nothing to do
		return nil
	}

	current, err := b.farosClientSet.TenancyV1alpha1().Workspaces().Get(ctx, w.Workspace.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	current.Spec = w.Workspace.Spec
	current.ResourceVersion = ""

	_, err = b.farosClientSet.TenancyV1alpha1().Workspaces().Update(ctx, &w.Workspace, metav1.UpdateOptions{})
	return err
}
