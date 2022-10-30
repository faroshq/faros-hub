package plugin

import (
	"context"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/kcp-dev/kcp/pkg/cliplugins/base"
	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
)

// GetWorkspacesOptions contains options for configuring faros workspaces
type GetWorkspacesOptions struct {
	*base.Options

	TenantWorkspaceAPI string
}

// NewGetWorkspacesOptions returns a new GetWorkspacesOptions.
func NewGetWorkspacesOptions(streams genericclioptions.IOStreams) *GetWorkspacesOptions {
	return &GetWorkspacesOptions{
		Options: base.NewOptions(streams),
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *GetWorkspacesOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)
}

// Complete ensures all dynamically populated fields are initialized.
func (o *GetWorkspacesOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	o.TenantWorkspaceAPI = "/apis/faros.sh/workspaces"

	return nil
}

// Validate validates the SyncOptions are complete and usable.
func (o *GetWorkspacesOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// Run gets workspaces from tenant workspace api
func (o *GetWorkspacesOptions) Run(ctx context.Context) error {
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

	workspaces := &tenancyv1alpha1.WorkspaceList{}

	err = farosclient.RESTClient().Get().AbsPath("/faros.sh/workspaces/bob").Do(ctx).Into(workspaces)
	if err != nil {
		return err
	}
	spew.Dump(workspaces)

	return err
}
