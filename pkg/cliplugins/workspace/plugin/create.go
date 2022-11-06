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

// GetWorkspacesOptions contains options for configuring faros workspaces
type CreateWorkspacesOptions struct {
	*base.Options

	TenantWorkspaceAPI string

	Name string

	Description string

	Members []string
}

// NewCreateWorkspacesOptions returns a new NewCreateWorkspacesOptions.
func NewCreateWorkspacesOptions(streams genericclioptions.IOStreams) *CreateWorkspacesOptions {
	return &CreateWorkspacesOptions{
		Options: base.NewOptions(streams),
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *CreateWorkspacesOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)

	cmd.Flags().StringArrayVarP(&o.Members, "members", "m", o.Members, "Additional members emails to add to the workspace")
	cmd.Flags().StringVarP(&o.Description, "description", "d", o.Description, "Description of the workspace")

}

// Complete ensures all dynamically populated fields are initialized.
func (o *CreateWorkspacesOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	o.TenantWorkspaceAPI = "/apis/faros.sh/workspaces"

	return nil
}

// Validate validates the WorkspacesOptions are complete and usable.
func (o *CreateWorkspacesOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	if o.Name == "" {
		errs = append(errs, fmt.Errorf("workspace name is required"))
	}

	return utilerrors.NewAggregate(errs)
}

// Run gets workspaces from tenant workspace api
func (o *CreateWorkspacesOptions) Run(ctx context.Context) error {
	config, err := o.ClientConfig.ClientConfig()
	if err != nil {
		return err
	}

	u, err := url.Parse(config.Host)
	if err != nil {
		return err
	}
	config.Host = u.Host

	farosclient, err := farosclient.NewForConfig(config)
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

	err = farosclient.RESTClient().Post().Body(patch).AbsPath("/faros.sh/workspaces").Do(ctx).Into(&workspace)
	if err != nil {
		return err
	}

	return nil
}
