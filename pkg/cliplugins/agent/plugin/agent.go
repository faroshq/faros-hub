package plugin

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"

	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
	"github.com/kcp-dev/kcp/pkg/cliplugins/helpers"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	edgevalpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	farosclient "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	"github.com/faroshq/faros-hub/pkg/cliplugins/base"
)

//go:embed *.yaml
var embeddedResources embed.FS

// GenerateOptions contains options for configuring a Agent and its corresponding process.
type GenerateOptions struct {
	*base.Options
	// OutputFile is the path to a file where the YAML for the agent kubeconfig should be written.
	OutputFile string
	// RegistrationName is name for Registration to use
	RegistrationName string
	// AgentName is the name of the Agent to be created.
	AgentName string
	// Namespace name
	Namespace string
}

// NewGenerateOptions returns a new GenerateOptions.
func NewGenerateOptions(streams genericclioptions.IOStreams) *GenerateOptions {
	return &GenerateOptions{
		Options: base.NewOptions(streams),
	}
}

// BindFlags binds fields GenerateOptions as command line flags to cmd's flagset.
func (o *GenerateOptions) BindFlags(cmd *cobra.Command) {
	o.Options.BindFlags(cmd)

	cmd.Flags().StringVarP(&o.OutputFile, "file", "f", o.OutputFile, "The manifest file to be created and applied to the physical cluster. Use - for stdout.")
	cmd.Flags().StringVarP(&o.RegistrationName, "registration", "r", o.RegistrationName, "Registration name to be used for agent.")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "default", "Namespace name")

}

// Complete ensures all dynamically populated fields are initialized.
func (o *GenerateOptions) Complete(args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	o.AgentName = args[0]

	return nil
}

// Validate validates the SyncOptions are complete and usable.
func (o *GenerateOptions) Validate() error {
	var errs []error

	if err := o.Options.Validate(); err != nil {
		errs = append(errs, err)
	}

	if o.OutputFile == "" {
		errs = append(errs, errors.New("--output-file is required"))
	}

	return utilerrors.NewAggregate(errs)
}

// Run prepares an agent kubeconfig for use with a agent and outputs the
// configuration required to deploy a agent to remote agent
func (o *GenerateOptions) Run(ctx context.Context) error {
	config, err := o.ClientConfig.ClientConfig()
	if err != nil {
		return err
	}

	var outputFile *os.File
	if o.OutputFile == "-" {
		outputFile = os.Stdout
	} else {
		outputFile, err = os.Create(o.OutputFile)
		if err != nil {
			return err
		}
		defer outputFile.Close()
	}

	token, err := o.enableAgentToRegister(ctx, config)
	if err != nil {
		return err
	}

	configURL, currentClusterName, err := helpers.ParseClusterURL(config.Host)
	if err != nil {
		return fmt.Errorf("current URL %q does not point to cluster workspace", config.Host)
	}

	// Compose the agent's upstream configuration server URL without any path. This is
	// required so long as the API importer and agent expect to require cluster clients.
	serverURL := configURL.Scheme + "://" + configURL.Host + "/clusters/" + currentClusterName.String()

	input := templateInput{
		AgentName:             o.AgentName,
		ServerURL:             serverURL,
		CAData:                base64.StdEncoding.EncodeToString(config.CAData),
		InsecureSkipTLSVerify: config.Insecure,
		Token:                 token,
		LogicalCluster:        currentClusterName.String(),
		Namespace:             o.Namespace,
	}

	resources, err := renderAgentResources(input)
	if err != nil {
		return err
	}

	_, err = outputFile.Write(resources)
	if o.OutputFile != "-" {
		fmt.Fprintf(o.Out, "\nWrote agent config to %s Use\n", o.OutputFile)
	}
	return err
}

// enableAgentToRegister gets individual kubeconfig for registration object
func (o *GenerateOptions) enableAgentToRegister(ctx context.Context, config *rest.Config) (saToken string, err error) {
	farosClient, err := farosclient.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create kcp client: %w", err)
	}
	coreClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create k8s client: %w", err)
	}

	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: o.Namespace,
		},
	}

	_, err = coreClient.CoreV1().Namespaces().Create(ctx, &namespace, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return "", fmt.Errorf("failed to create namespace: %w", err)
	}

	registrationName := o.RegistrationName
	if registrationName == "" {
		registrationName = o.AgentName
	}

	template := edgevalpha1.Registration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      registrationName,
			Namespace: o.Namespace,
		},
	}

	registration, err := farosClient.EdgeV1alpha1().Registrations(o.Namespace).Get(ctx, template.Name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		fmt.Fprintf(o.Out, "Creating registration %s\n", registration.Name)
		registration, err = farosClient.EdgeV1alpha1().Registrations(o.Namespace).Create(ctx, &template, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return "", fmt.Errorf("failed to create registration %q: %w", registration.Name, err)
		}
	case err == nil:
		return registration.Status.Token, nil
	default:
		return "", fmt.Errorf("failed to create the ClusterRole %s", err)
	}

	// wait for registration to be ready
	fmt.Fprintf(o.Out, "Waiting for registration %s to be ready\n", registration.Name)
	err = wait.PollImmediateWithContext(ctx, 100*time.Millisecond, 20*time.Second, func(ctx context.Context) (bool, error) {
		registration, err = farosClient.EdgeV1alpha1().Registrations(o.Namespace).Get(ctx, registration.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(o.ErrOut, "failed to retrieve Registration: %v", err)
			return false, nil
		}
		return conditions.IsTrue(registration, conditionsv1alpha1.ReadyCondition) && registration.Status.Token != "", nil
	})
	return registration.Status.Token, err

}

// templateInput represents the external input required to render the resources to
// deploy the syncer to a pcluster.
type templateInput struct {
	// AgentName is the name of the agent
	AgentName string
	// ServerURL is the logical cluster url the syncer configuration will use
	ServerURL string
	// CAData holds the PEM-encoded bytes of the ca certificate(s) a syncer will use to validate
	// kcp's serving certificate
	CAData string
	// InsecureSkipTLSVerify controls whether a syncer verifies the server's certificate chain and host name
	InsecureSkipTLSVerify bool
	// Token is the service account token used to authenticate a syncer for access to a workspace
	Token string
	// Namespace is the name of the syncer namespace on the pcluster
	Namespace string
	// LogicalCluster is the qualified kcp logical cluster name the syncer will sync from
	LogicalCluster string
}

// templateArgs represents the full set of arguments required to render the resources
// required to deploy the syncer.
type templateArgs struct {
	templateInput
}

// renderAgentResources renders the resources required to deploy an agent
func renderAgentResources(input templateInput) ([]byte, error) {
	tmplArgs := templateArgs{
		templateInput: input,
	}

	agentTemplate, err := embeddedResources.ReadFile("agent.yaml")
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("agentTemplate").Parse(string(agentTemplate))
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buffer, tmplArgs)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
