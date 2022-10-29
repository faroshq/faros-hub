package plugin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/faroshq/faros-hub/pkg/models"
	"github.com/kcp-dev/kcp/pkg/cliplugins/base"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog"
	"sigs.k8s.io/yaml"
)

// LoginSetupOptions contains options for login via faros API
type LoginSetupOptions struct {
	*base.Options

	// ConfigFile of CLI config
	ConfigFile string
}

// NewGenerateOptions returns a new GenerateOptions.
func NewLoginSetupOptions(streams genericclioptions.IOStreams) *LoginSetupOptions {
	return &LoginSetupOptions{
		Options: base.NewOptions(streams),
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *LoginSetupOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)

	homedir, err := os.UserHomeDir()
	if err != nil {
		klog.Error("Failed to get user home directory")
		homedir = "/tmp/"
	}

	cmd.Flags().StringVarP(&o.ConfigFile, "config", "c", filepath.Join(homedir, ".faros/config.yaml"), "Faros CLI config location")
}

// Complete ensures all dynamically populated fields are initialized.
func (o *LoginSetupOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
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
	fmt.Println("Logging into Faros Hub...")

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
		return o.configureCLI(ctx, *response)
	case err := <-errCh:
		return fmt.Errorf("trying to authorize the client: %s", err)

	case <-time.After(60 * time.Second):
		return errors.New("authorization request timed out waiting to complete")
	}

}

func (o *LoginSetupOptions) configureCLI(ctx context.Context, response models.LoginResponse) error {
	fmt.Printf("Persisting login configuration to %s \n", o.ConfigFile)
	data, err := os.ReadFile(o.ConfigFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	base := filepath.Dir(o.ConfigFile)
	err = os.MkdirAll(base, 0755)
	if err != nil {
		return err
	}

	config := models.NewCLIConfig()
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	config.Spec.Token = response.RawIDToken
	config.Spec.BaseURL = response.ServerBaseURL
	config.Spec.Email = response.Email

	data, err = yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(o.ConfigFile, data, 0644)
}
