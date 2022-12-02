package plugin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/cliplugins/base"
)

var kubeConfigAuthKey = "faros"

// EnableOptions contains options for configuring faros plugins
type EnableOptions struct {
	*base.Options

	Name       string
	PluginName string
	Version    string
	Workspace  string
	Namespace  string
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

	cmd.Flags().StringVar(&o.Version, "plugin-version", "latest", "Plugin version to use")
	cmd.Flags().StringVar(&o.PluginName, "plugin-name", "", "Plugin name")
	cmd.Flags().StringVar(&o.Workspace, "workspace", "default", "Workspace name")
}

// Complete ensures all dynamically populated fields are initialized.
func (o *EnableOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	if o.PluginName == "" {
		o.PluginName = o.Name // default to request name
	}

	if o.Version == "" {
		o.Version = "latest"
	}

	rawConfig, err := o.ClientConfig.RawConfig()
	if err != nil {
		fmt.Printf("Not able to determine workspace name from kubeconfig: %v", err)
		return err
	}
	if rawConfig.CurrentContext == kubeConfigAuthKey {
		return fmt.Errorf("kubeconfig not set to any workspace. Use `faros workspace use` to set the current workspace")
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

	farosclient, err := farosclient.NewForConfig(config)
	if err != nil {
		return err
	}

	plugin := pluginsv1alpha1.Request{
		TypeMeta: metav1.TypeMeta{
			Kind:       pluginsv1alpha1.RequestKind,
			APIVersion: pluginsv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: o.Name,
		},
		Spec: pluginsv1alpha1.RequestSpec{
			Version: o.Version,
			Name:    o.PluginName,
		},
	}

	_, err = farosclient.PluginsV1alpha1().Requests().Create(ctx, &plugin, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "Plugin request %s created", o.Name)
	return nil
}
