package bridge

import (
	"context"

	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/faroshq/faros-hub/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (b *bridge) createUser(ctx context.Context, userID string) error {
	u, err := b.store.GetUser(ctx, models.User{
		ID: userID,
	})
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}
	if err == store.ErrRecordNotFound { // gone, nothing to do
		return nil
	}

	user := u.User
	if user.Labels == nil {
		user.Labels = map[string]string{}
	}
	user.Labels[models.LabelIdentifier] = userID

	_, err = b.farosClientSet.TenancyV1alpha1().Users().Create(ctx, &user, metav1.CreateOptions{})
	return err
}

func (b *bridge) deleteUser(ctx context.Context, userID string) error {
	err := b.farosClientSet.TenancyV1alpha1().Users().DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: models.LabelIdentifier + "=" + userID,
	})
	return err
}

func (b *bridge) updateUser(ctx context.Context, userID string) error {
	w, err := b.store.GetUser(ctx, models.User{
		ID: userID,
	})
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}
	if err == store.ErrRecordNotFound { // gone, nothing to do
		return nil
	}

	current, err := b.farosClientSet.TenancyV1alpha1().Users().Get(ctx, w.User.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	current.Spec = w.User.Spec
	current.ResourceVersion = ""

	_, err = b.farosClientSet.TenancyV1alpha1().Users().Update(ctx, &w.User, metav1.UpdateOptions{})
	return err
}
