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

// EnableOptions contains options for configuring faros plugins
type EnableOptions struct {
	*base.Options
	Name string
}

// NewEnableOptions returns a new EnableOptions.
func NewEnableOptions(streams genericclioptions.IOStreams) *EnableOptions {
	return &EnableOptions{
		Options: base.NewOptions(streams),
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *EnableOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)
}

// Complete ensures all dynamically populated fields are initialized.
func (o *EnableOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	return nil
}

// Validate validates the SyncOptions are complete and usable.
func (o *EnableOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// Run gets plugins from tenant workspace api
func (o *EnableOptions) Run(ctx context.Context) error {
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
