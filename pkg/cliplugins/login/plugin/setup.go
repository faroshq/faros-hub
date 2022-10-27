package plugin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/faroshq/faros-hub/pkg/models"
	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/cliplugins/base"
	pluginhelpers "github.com/kcp-dev/kcp/pkg/cliplugins/helpers"
	"github.com/kcp-dev/logicalcluster/v2"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// LoginSetupOptions contains options for login via faros API
type LoginSetupOptions struct {
	*base.Options

	startingConfig *clientcmdapi.Config

	// Name is the name of the workspace to switch to.
	Name string
	// Overwrite indicates the context should be updated if it already exists. This is required to perform the update.
	Overwrite bool

	// for testing
	modifyConfig func(configAccess clientcmd.ConfigAccess, newConfig *clientcmdapi.Config) error
}

// NewGenerateOptions returns a new GenerateOptions.
func NewLoginSetupOptions(streams genericclioptions.IOStreams) *LoginSetupOptions {
	return &LoginSetupOptions{
		Options: base.NewOptions(streams),
		modifyConfig: func(configAccess clientcmd.ConfigAccess, newConfig *clientcmdapi.Config) error {
			return clientcmd.ModifyConfig(configAccess, *newConfig, true)
		},
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *LoginSetupOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)

	cmd.Flags().BoolVar(&o.Overwrite, "overwrite", o.Overwrite, "Overwrite the context if it already exists")

}

// Complete ensures all dynamically populated fields are initialized.
func (o *LoginSetupOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	var err error
	o.startingConfig, err = o.ClientConfig.ConfigAccess().GetStartingConfig()
	if err != nil {
		return err
	}

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	return nil
}

// Validate validates the inputs
func (o *LoginSetupOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// Run prepares initiated login flow via IDP
func (o *LoginSetupOptions) Run(ctx context.Context) error {
	fmt.Println("running login setup")
	// TODO: URL here should be configurable to
	// faros server so cookies are preserved
	//rest, err := o.ClientConfig.ClientConfig()
	//////if err != nil {
	//////	return err
	//}

	//u, err := url.Parse(rest.Host)
	//////if err != nil {
	//////	return err
	//}

	doneCh := make(chan struct{})
	errCh := make(chan error)
	response := &models.LoginResponse{}

	// local server to catch the response
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		result, err := handleLoginCallback(req, w)
		if err != nil {
			errCh <- err
			return
		}
		response = result
		doneCh <- struct{}{}
	})

	l, err := getLocalListener()
	if err != nil {
		return err
	}

	// start serving locally and wait for the response
	go func() {
		if err := http.Serve(l, http.DefaultServeMux); err != nil {
			errCh <- fmt.Errorf("trying to start local http server: %s", err)
		}
	}()

	url := fmt.Sprintf("https://localhost:6443/faros.sh/oidc/login?redirect_uri=http://localhost:%d", l.Addr().(*net.TCPAddr).Port)
	spew.Dump(url)

	if err := open.Run(url); err != nil {
		return fmt.Errorf("trying to open web browser, error: %s", err)
	}

	// wait for the response
	select {
	case <-doneCh:
		return o.configureKubeconfig(ctx, *response)
	case err := <-errCh:
		return fmt.Errorf("trying to authorize the client: %s", err)

	case <-time.After(60 * time.Second):
		return errors.New("authorization request timed out waiting to complete")
	}

}

func (o *LoginSetupOptions) configureKubeconfig(ctx context.Context, response models.LoginResponse) error {
	config, err := o.ClientConfig.RawConfig()
	if err != nil {
		return err
	}
	currentContext, ok := config.Contexts[config.CurrentContext]
	if !ok {
		return fmt.Errorf("current context %q is not found in kubeconfig", config.CurrentContext)
	}
	currentCluster, ok := config.Clusters[currentContext.Cluster]
	if !ok {
		return fmt.Errorf("current cluster %q is not found in kubeconfig", currentContext.Cluster)
	}
	_, currentClusterName, err := pluginhelpers.ParseClusterURL(currentCluster.Server)
	if err != nil {
		return fmt.Errorf("current URL %q does not point to cluster workspace", currentCluster.Server)
	}

	newServerHost := currentCluster.Server
	if currentClusterName.String() != o.Name {
		cluster := logicalcluster.New(o.Name)
		if strings.Contains(o.Name, ":") && !cluster.HasPrefix(logicalcluster.New("system")) &&
			!cluster.HasPrefix(tenancyv1alpha1.RootCluster) {
			return fmt.Errorf("invalid workspace name format: %s", o.Name)
		}

		config, err := o.ClientConfig.ClientConfig()
		if err != nil {
			return err
		}

		u, _, err := pluginhelpers.ParseClusterURL(config.Host)
		if err != nil {
			return fmt.Errorf("current URL %q does not point to cluster workspace", config.Host)
		}

		u.Path = path.Join(u.Path, cluster.Path())
		newServerHost = u.String()
	}

	if o.Name == "" {
		o.Name = currentClusterName.String()
	}

	_, existedBefore := o.startingConfig.Contexts[o.Name]
	if existedBefore && !o.Overwrite {
		return fmt.Errorf("context %q already exists in kubeconfig, use --overwrite to update it", o.Name)
	}

	newKubeConfig := o.startingConfig.DeepCopy()
	newCluster := *currentCluster
	newCluster.Server = newServerHost
	newKubeConfig.Clusters[o.Name] = &newCluster
	newContext := *currentContext
	newContext.Cluster = o.Name
	newContext.AuthInfo = response.Email
	newKubeConfig.Contexts[o.Name] = &newContext
	newKubeConfig.CurrentContext = o.Name
	newKubeConfig.AuthInfos[response.Email] = &clientcmdapi.AuthInfo{
		Token: response.RawIDToken,
	}

	if err := o.modifyConfig(o.ClientConfig.ConfigAccess(), newKubeConfig); err != nil {
		return err
	}

	if existedBefore {
		if o.startingConfig.CurrentContext == o.Name {
			_, err = fmt.Fprintf(o.Out, "Updated context %q.\n", o.Name)
		} else {
			_, err = fmt.Fprintf(o.Out, "Updated context %q and switched to it.\n", o.Name)
		}
	} else {
		_, err = fmt.Fprintf(o.Out, "Created context %q and switched to it.\n", o.Name)
	}

	return nil
}
