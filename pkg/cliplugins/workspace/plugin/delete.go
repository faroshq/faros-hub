package plugin

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/cliplugins/base"
)

// DeleteWorkspacesOptions contains options for configuring faros workspaces
type DeleteWorkspacesOptions struct {
	*base.Options

	Name string
}

// NewGetWorkspacesOptions returns a new GetWorkspacesOptions.
func NewDeleteWorkspacesOptions(streams genericclioptions.IOStreams) *DeleteWorkspacesOptions {
	return &DeleteWorkspacesOptions{
		Options: base.NewOptions(streams),
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *DeleteWorkspacesOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)
}

// Complete ensures all dynamically populated fields are initialized.
func (o *DeleteWorkspacesOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	return nil
}

// Validate validates the SyncOptions are complete and usable.
func (o *DeleteWorkspacesOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// Run gets workspaces from tenant workspace api
func (o *DeleteWorkspacesOptions) Run(ctx context.Context) error {
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

	workspace := &tenancyv1alpha1.Workspace{}

	err = farosclient.RESTClient().Delete().AbsPath("/faros.sh/workspaces/" + o.Name).Do(ctx).Into(workspace)
	if err != nil {
		return err
	}

	fmt.Println("Workspace deleted successfully")
	return nil
}
