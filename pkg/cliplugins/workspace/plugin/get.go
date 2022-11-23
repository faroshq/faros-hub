package plugin

import (
	"context"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/cliplugins/base"
	utilprint "github.com/faroshq/faros-hub/pkg/util/print"
)

// GetOptions contains options for configuring faros
type GetOptions struct {
	*base.Options
	Name string
}

// NewGetOptions returns a new GetOptions.
func NewGetOptions(streams genericclioptions.IOStreams) *GetOptions {
	return &GetOptions{
		Options: base.NewOptions(streams),
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *GetOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)
}

// Complete ensures all dynamically populated fields are initialized.
func (o *GetOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	return nil
}

// Validate validates the SyncOptions are complete and usable.
func (o *GetOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// Run gets workspaces from tenant workspace api
func (o *GetOptions) Run(ctx context.Context) error {
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

	err = farosclient.RESTClient().Get().AbsPath("/faros.sh/api/v1alpha1/workspaces").Do(ctx).Into(workspaces)
	if err != nil {
		return err
	}

	// drop managed fields
	for i := range workspaces.Items {
		workspaces.Items[i].ObjectMeta.ManagedFields = nil
	}

	if o.Output == utilprint.FormatTable {
		table := utilprint.DefaultTable()
		table.SetHeader([]string{"NAME", "MEMBERS", "DESCRIPTION", "STATUS", "AGE"})
		for _, workspace := range workspaces.Items {
			{
				status := "Unknown"
				if len(workspace.Status.Conditions) > 0 {
					status = string(workspace.Status.Conditions[0].Type)
				}

				table.Append([]string{
					workspace.Name,
					strings.Join(workspace.Spec.Members, ","),
					workspace.Spec.Description,
					status,
					utilprint.Since(workspace.CreationTimestamp.Time).String()},
				)
			}
		}
		table.Render()
		return nil
	}

	return utilprint.PrintWithFormat(workspaces, o.Output)
}
