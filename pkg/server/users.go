package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/kcp-dev/logicalcluster/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
)

const (
	// UserLabel is the label used to identify the user
	UserLabel = "faros.sh/user"
)

// registerOrUpdateUser will register or update user in the system when user is authenticated
func (s *Service) registerOrUpdateUser(ctx context.Context, user *tenancyv1alpha1.User) (*tenancyv1alpha1.User, error) {
	if user.Name == "" {
		// Hack to it does not start with a number
		// https://github.com/kcp-dev/kcp/blob/main/pkg/server/filters/filters.go#L52
		for {
			user.Name = string(uuid.NewUUID())
			if user.Name[0] > '9' {
				break
			}
		}
	}

	// we will be selecting based on labels, but k8s does not allow symbols like '@' in labels
	// so we will replace it with '-at-' checking before if only single @ exists

	if !strings.Contains(user.Spec.Email, "@") {
		return nil, fmt.Errorf("invalid email address")
	}

	labelEmail := strings.Replace(user.Spec.Email, "@", "-at-", 1)

	cluster := logicalcluster.New(s.config.ControllersTenantWorkspace)

	users, err := s.farosClient.Cluster(cluster).TenancyV1alpha1().Users().List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", UserLabel, labelEmail),
	})
	if err != nil {
		return nil, err
	}
	// existing user flow
	if len(users.Items) > 0 {
		if len(users.Items) > 1 {
			return nil, fmt.Errorf("multiple users found with same email")
		}

		current := &users.Items[0]
		current.Spec = user.Spec
		if current.Labels == nil {
			current.Labels = make(map[string]string)
		}
		current.Labels[UserLabel] = labelEmail
		return s.farosClient.Cluster(cluster).TenancyV1alpha1().Users().Update(ctx, current, metav1.UpdateOptions{})
	} else {
		user.Labels = map[string]string{
			UserLabel: labelEmail,
		}
		return s.farosClient.Cluster(cluster).TenancyV1alpha1().Users().Create(ctx, user, metav1.CreateOptions{})
	}
}
