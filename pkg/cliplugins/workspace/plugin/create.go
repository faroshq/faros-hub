package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/cliplugins/base"
)

// GetOptions contains options for configuring faros workspaces
type CreateOptions struct {
	*base.Options
	Name        string
	Description string
	Members     []string
}

// NewCreateOptions returns a new NewCreateOptions.
func NewCreateOptions(streams genericclioptions.IOStreams) *CreateOptions {
	return &CreateOptions{
		Options: base.NewOptions(streams),
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *CreateOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)

	cmd.Flags().StringArrayVarP(&o.Members, "members", "m", o.Members, "Additional members emails to add to the workspace")
	cmd.Flags().StringVarP(&o.Description, "description", "d", o.Description, "Description of the workspace")

}

// Complete ensures all dynamically populated fields are initialized.
func (o *CreateOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	return nil
}

// Validate validates the Options are complete and usable.
func (o *CreateOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	if o.Name == "" {
		errs = append(errs, fmt.Errorf("workspace name is required"))
	}

	return utilerrors.NewAggregate(errs)
}

// Run gets  from tenant workspace api
func (o *CreateOptions) Run(ctx context.Context) error {
	config, err := o.ClientConfig.ClientConfig()
	if err != nil {
		return err
	}

	u, err := url.Parse(config.Host)
	if err != nil {
		return err
	}
	config.Host = u.Host

	farosClient, err := farosclient.NewForConfig(config)
	if err != nil {
		return err
	}

	workspace := tenancyv1alpha1.Workspace{
		TypeMeta: metav1.TypeMeta{
			Kind:       tenancyv1alpha1.WorkspaceKind,
			APIVersion: tenancyv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: o.Name,
		},
		Spec: tenancyv1alpha1.WorkspaceSpec{
			Description: o.Description,
			Members:     o.Members,
		},
	}

	patch, err := json.Marshal(workspace)
	if err != nil {
		return fmt.Errorf("error creating patch: %v", err)
	}

	err = farosClient.RESTClient().Post().Body(patch).AbsPath("/faros.sh/api/v1alpha1/workspaces").Do(ctx).Into(&workspace)
	if err != nil {
		return err
	}

	return nil
}
