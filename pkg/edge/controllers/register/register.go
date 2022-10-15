package register

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
)

type Register interface {
	// Register registers a new agent
	Register(ctx context.Context, name string, namespace string) error
}

type register struct {
	config *rest.Config
	client farosclient.Interface
}

func New(config *rest.Config) (Register, error) {
	farosclient, err := farosclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create faros client: %w", err)
	}
	return &register{
		config: config,
		client: farosclient,
	}, nil
}

func (r *register) Register(ctx context.Context, name, namespace string) error {
	_, err := r.client.EdgeV1alpha1().Agents(namespace).Create(ctx, &edgev1alpha1.Agent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to register agent: %w", err)
	}
	return nil
}
