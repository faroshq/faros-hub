package plugin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/kcp-dev/kcp/pkg/cliplugins/base"
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

	if o.Name == "" && len(args) > 0 {
		o.Name = args[0]
	}

	return nil
}

// Validate validates the inputs
func (o *LoginSetupOptions) Validate() error {
	var errs []error

	if o.Name == "" {
		errs = append(errs, errors.New("workspace name is required"))
	}

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

// Run prepares initiated login flow via IDP
func (o *LoginSetupOptions) Run(ctx context.Context) error {
	fmt.Println("running login setup")

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

	if err := open.Run(url); err != nil {
		return fmt.Errorf("trying to open web browser, error: %s", err)
	}

	// wait for the response
	select {
	case <-doneCh:
		return o.configureKubeConfig(ctx, *response)
	case err := <-errCh:
		return fmt.Errorf("trying to authorize the client: %s", err)

	case <-time.After(60 * time.Second):
		return errors.New("authorization request timed out waiting to complete")
	}

}

func (o *LoginSetupOptions) configureKubeConfig(ctx context.Context, response models.LoginResponse) error {
	config, err := o.ClientConfig.RawConfig()
	if err != nil {
		return err
	}

	// setup cluster
	cluster, exists := config.Clusters[o.Name]
	if !exists {
		cluster = clientcmdapi.NewCluster()
	}
	cluster.Server = response.ServerBaseURL + "/" + o.Name
	if response.CertificateAuthorityData == "" {
		cluster.InsecureSkipTLSVerify = true
	} else {
		cluster.CertificateAuthorityData = []byte(response.CertificateAuthorityData)
	}
	config.Clusters[o.Name] = cluster

	// setup user
	user, exists := config.AuthInfos[o.Name]
	if !exists {
		user = clientcmdapi.NewAuthInfo()
	}
	user.Token = response.RawIDToken
	config.AuthInfos[o.Name] = user

	// setup context
	context, exists := config.Contexts[o.Name]
	if !exists {
		context = clientcmdapi.NewContext()
	}
	context.Cluster = o.Name
	context.AuthInfo = o.Name
	config.Contexts[o.Name] = context

	config.CurrentContext = o.Name

	return o.modifyConfig(o.ClientConfig.ConfigAccess(), &config)
}
