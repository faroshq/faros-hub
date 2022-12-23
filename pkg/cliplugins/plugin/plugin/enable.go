package plugin

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/cliplugins/base"
)

var kubeConfigAuthKey = "faros"

var pluginNameRegexp = regexp.MustCompile(`^[a-z]{0,10}.[a-z]{0,10}.plugins.faros.sh?$`)

// EnableOptions contains options for configuring faros plugins
type EnableOptions struct {
	*base.Options

	Name       string
	PluginName string
	Version    string
	Workspace  string
	Namespace  string

	BindingLabelSelector         map[string]string
	bindingLabelSelectorInternal []string
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
	cmd.Flags().StringVar(&o.Namespace, "namespace", "default", "Namespace name")
	cmd.Flags().StringSliceVar(&o.bindingLabelSelectorInternal, "binding-label-selector", []string{}, "Label selector for binding")
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
		o.PluginName = o.Name // default to name
	}

	if o.Version == "" {
		o.Version = "latest"
	}

	if o.BindingLabelSelector == nil {
		o.BindingLabelSelector = make(map[string]string)
	}
	for _, s := range o.bindingLabelSelectorInternal {

		if !strings.Contains(s, "=") {
			return fmt.Errorf("invalid label selector. Example: location=vilnius,env=dev: %s", s)
		}
		parts := strings.Split(s, "=")
		o.BindingLabelSelector[parts[0]] = parts[1]
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

	if !pluginNameRegexp.MatchString(o.PluginName) {
		errs = append(errs, fmt.Errorf("invalid plugin name: %s", o.PluginName))
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

	pluginRequest := pluginsv1alpha1.Request{
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

	_, err = farosclient.PluginsV1alpha1().Requests().Create(ctx, &pluginRequest, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// wait for plugin to be ready
	fmt.Fprintf(o.Out, "Waiting for plugin request %s to be ready\n", pluginRequest.Name)
	err = wait.PollImmediateWithContext(ctx, 100*time.Millisecond, 60*time.Second, func(ctx context.Context) (bool, error) {
		pluginRequestCurrent, err := farosclient.PluginsV1alpha1().Requests().Get(ctx, pluginRequest.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(o.ErrOut, "failed to retrieve Registration: %v", err)
			return false, nil
		}
		return conditions.IsTrue(pluginRequestCurrent, conditionsv1alpha1.ReadyCondition), nil
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "Plugin request %s created\n", o.Name)
	if len(o.BindingLabelSelector) > 0 {
		fmt.Fprintf(o.Out, "Creating plugin binding\n")
		// found the api version
		binding := pluginsv1alpha1.Binding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      o.Name,
				Namespace: o.Namespace,
			},
			Spec: pluginsv1alpha1.BindingSpec{
				PluginType: o.PluginName,
				PluginName: o.Name,
				Selector: metav1.LabelSelector{
					MatchLabels: o.BindingLabelSelector,
				},
			},
		}

		_, err = farosclient.PluginsV1alpha1().Bindings(o.Namespace).Create(ctx, &binding, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		fmt.Fprintf(o.Out, "Plugin binding %s created\n", o.Name)
		return nil
	}

	return nil
}
