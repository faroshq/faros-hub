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

// UseWorkspacesOptions contains options for configuring faros workspaces
type UseWorkspacesOptions struct {
	*base.Options
	Name string

	// for testing
	modifyConfig func(configAccess clientcmd.ConfigAccess, newConfig *clientcmdapi.Config) error
}

// NewUseWorkspacesOptions returns a new GetWorkspacesOptions.
func NewUseWorkspacesOptions(streams genericclioptions.IOStreams) *UseWorkspacesOptions {
	return &UseWorkspacesOptions{
		Options: base.NewOptions(streams),
		modifyConfig: func(configAccess clientcmd.ConfigAccess, newConfig *clientcmdapi.Config) error {
			return clientcmd.ModifyConfig(configAccess, *newConfig, true)
		},
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *UseWorkspacesOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)
}

// Complete ensures all dynamically populated fields are initialized.
func (o *UseWorkspacesOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	return nil
}

// Validate validates the SyncOptions are complete and usable.
func (o *UseWorkspacesOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// Run gets workspaces from tenant workspace api
func (o *UseWorkspacesOptions) Run(ctx context.Context) error {
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

	err = farosclient.RESTClient().Get().AbsPath("/faros.sh/api/v1alpha1/workspaces/" + o.Name).Do(ctx).Into(workspace)
	if err != nil {
		return err
	}

	// Get raw config and add new cluster and context to it
	rawConfig, err := o.ClientConfig.RawConfig()
	if err != nil {
		return err
	}

	rawConfig.Clusters[workspace.Name] = &clientcmdapi.Cluster{
		Server: workspace.Status.WorkspaceURL,
	}

	farosCluster, ok := rawConfig.Clusters[kubeConfigAuthKey]
	if !ok {
		rawConfig.Clusters[workspace.Name].InsecureSkipTLSVerify = true
	} else {
		if farosCluster.InsecureSkipTLSVerify {
			rawConfig.Clusters[workspace.Name].InsecureSkipTLSVerify = true
		} else {
			rawConfig.Clusters[workspace.Name].CertificateAuthorityData = farosCluster.CertificateAuthorityData
			rawConfig.Clusters[workspace.Name].CertificateAuthority = farosCluster.CertificateAuthority
		}
	}

	rawConfig.Contexts[workspace.Name] = &clientcmdapi.Context{
		Cluster:  workspace.Name,
		AuthInfo: kubeConfigAuthKey,
	}

	rawConfig.CurrentContext = workspace.Name

	fmt.Println("Using workspace", workspace.Name)
	return o.modifyConfig(o.ClientConfig.ConfigAccess(), &rawConfig)
}
