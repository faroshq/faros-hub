package cmd

import (
	"fmt"

	"github.com/faroshq/faros-hub/pkg/cliplugins/login/plugin"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	loginExample = `
	# Login to faros with SSO
	%[1]s
`
)

// New provides a cobra command for workload operations.
func New(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Aliases:          []string{"login"},
		Use:              "login",
		Short:            "Manages Faros access via SSO",
		SilenceUsage:     true,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	generateOptions := plugin.NewLoginSetupOptions(streams)

	loginSetupCmd := &cobra.Command{
		Use:          "setup",
		Short:        "Setup initial login",
		Example:      fmt.Sprintf(loginExample, "--example-flags"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := generateOptions.Complete(args); err != nil {
				return err
			}

			if err := generateOptions.Validate(); err != nil {
				return err
			}

			return generateOptions.Run(c.Context())
		},
	}

	generateOptions.BindFlags(loginSetupCmd)
	cmd.AddCommand(loginSetupCmd)

	return cmd, nil
}
