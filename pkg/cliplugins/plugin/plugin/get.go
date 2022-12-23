package plugin

import (
	"context"
	"net/url"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/cliplugins/base"
	utilprint "github.com/faroshq/faros-hub/pkg/util/print"
)

// GetOptions contains options for configuring faros plugins
type GetOptions struct {
	*base.Options
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

// Run gets plugins from tenant workspace api
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

	plugins := &pluginsv1alpha1.PluginList{}

	// If context is set to 'faros', this means we are at root faros context.
	// If user overrides it or its set to something else, we assume we are in a workspace context.
	err = farosclient.RESTClient().Get().AbsPath("/faros.sh/api/v1alpha1/plugins").Do(ctx).Into(plugins)
	if err != nil {
		return err
	}

	// drop managed fields
	for i := range plugins.Items {
		plugins.Items[i].ObjectMeta.ManagedFields = nil
	}

	if o.Output == utilprint.FormatTable {
		table := utilprint.DefaultTable()
		table.SetHeader([]string{"NAME", "VERSION", "DESCRIPTION"})
		for _, plugin := range plugins.Items {
			{
				table.Append([]string{
					plugin.Name,
					plugin.Spec.Version,
					plugin.Spec.Description,
				})
			}
		}
		table.Render()
		return nil
	}

	return utilprint.PrintWithFormat(plugins, o.Output)
}
