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
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var kubeConfigAuthKey = "faros"

// UseOptions contains options for configuring faros
type UseOptions struct {
	*base.Options
	Name string

	// for testing
	modifyConfig func(configAccess clientcmd.ConfigAccess, newConfig *clientcmdapi.Config) error
}

// NewUseOptions returns a new GetOptions.
func NewUseOptions(streams genericclioptions.IOStreams) *UseOptions {
	return &UseOptions{
		Options: base.NewOptions(streams),
		modifyConfig: func(configAccess clientcmd.ConfigAccess, newConfig *clientcmdapi.Config) error {
			return clientcmd.ModifyConfig(configAccess, *newConfig, true)
		},
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *UseOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)
}

// Complete ensures all dynamically populated fields are initialized.
func (o *UseOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	return nil
}

// Validate validates the SyncOptions are complete and usable.
func (o *UseOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// Run gets workspace from tenant workspace api
func (o *UseOptions) Run(ctx context.Context) error {
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

	// Get raw config and add new cluster and context to it
	rawConfig, err := o.ClientConfig.RawConfig()
	if err != nil {
		return err
	}

	if o.Name != kubeConfigAuthKey {
		err = farosclient.RESTClient().Get().AbsPath("/faros.sh/api/v1alpha1/workspaces/" + o.Name).Do(ctx).Into(workspace)
		if err != nil {
			return err
		}

		rawConfig.Clusters[workspace.Spec.Name] = &clientcmdapi.Cluster{
			Server: workspace.Status.WorkspaceURL,
		}

		farosCluster, ok := rawConfig.Clusters[kubeConfigAuthKey]
		if !ok {
			rawConfig.Clusters[workspace.Spec.Name].InsecureSkipTLSVerify = true
		} else {
			if farosCluster.InsecureSkipTLSVerify {
				rawConfig.Clusters[workspace.Spec.Name].InsecureSkipTLSVerify = true
			} else {
				rawConfig.Clusters[workspace.Spec.Name].CertificateAuthorityData = farosCluster.CertificateAuthorityData
				rawConfig.Clusters[workspace.Spec.Name].CertificateAuthority = farosCluster.CertificateAuthority
			}
		}

		rawConfig.Contexts[workspace.Spec.Name] = &clientcmdapi.Context{
			Cluster:  workspace.Spec.Name,
			AuthInfo: kubeConfigAuthKey,
		}

		rawConfig.CurrentContext = workspace.Spec.Name

	} else {
		// if user requests "faros" context, just set it as current context
		rawConfig.CurrentContext = kubeConfigAuthKey
	}

	fmt.Println("Using workspace", o.Name)
	return o.modifyConfig(o.ClientConfig.ConfigAccess(), &rawConfig)
}
