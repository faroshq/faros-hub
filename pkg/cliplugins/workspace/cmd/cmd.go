package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/faroshq/faros-hub/pkg/cliplugins/workspace/plugin"
)

// New provides a cobra command for workload operations.
func New(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Aliases:          []string{"workspaces"},
		Use:              "workspace",
		Short:            "Manages workspaces",
		SilenceUsage:     true,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	getWorkspacesOptions := plugin.NewGetWorkspacesOptions(streams)

	getWorkspacesCmd := &cobra.Command{
		Use:          "get",
		Short:        "Get a workspaces",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := getWorkspacesOptions.Complete(args); err != nil {
				return err
			}

			if err := getWorkspacesOptions.Validate(); err != nil {
				return err
			}

			return getWorkspacesOptions.Run(c.Context())
		},
	}

	getWorkspacesOptions.BindFlags(getWorkspacesCmd)
	cmd.AddCommand(getWorkspacesCmd)

	return cmd, nil
}
