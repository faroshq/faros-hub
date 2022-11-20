package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/faroshq/faros-hub/pkg/cliplugins/agent/plugin"
)

var (
	agentExample = `
	# Ensure a agent is running on the specified agebt target.
	%[1]s <registration-name> -o agent.kubeconfig
	KUBECONFIG=<pcluster-config> <agent_image>
`
)

// New provides a cobra command for workload operations.
func New(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Aliases:          []string{"agents"},
		Use:              "agent",
		Short:            "Manages cluster agent targets",
		SilenceUsage:     true,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Sync command
	generateOptions := plugin.NewGenerateOptions(streams)

	generateAgentCmd := &cobra.Command{
		Use:          "generate <agent-name> --registration <registration-name> -f <output-file>",
		Short:        "Create a agent config in faros to be deployed in agent and callback into faros",
		Example:      fmt.Sprintf(agentExample, "kubectl faros agent generate"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return c.Help()
			}

			if err := generateOptions.Complete(args); err != nil {
				return err
			}

			if err := generateOptions.Validate(); err != nil {
				return err
			}

			return generateOptions.Run(c.Context())
		},
	}

	generateOptions.BindFlags(generateAgentCmd)
	cmd.AddCommand(generateAgentCmd)

	return cmd, nil
}
