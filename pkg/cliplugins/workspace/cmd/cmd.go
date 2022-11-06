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

	createWorkspacesOptions := plugin.NewCreateWorkspacesOptions(streams)
	createWorkspacesCmd := &cobra.Command{
		Use:          "create",
		Short:        "Create a workspaces",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := createWorkspacesOptions.Complete(args); err != nil {
				return err
			}

			if err := createWorkspacesOptions.Validate(); err != nil {
				return err
			}

			return createWorkspacesOptions.Run(c.Context())
		},
	}

	deleteWorkspacesOptions := plugin.NewDeleteWorkspacesOptions(streams)
	deleteWorkspacesCmd := &cobra.Command{
		Use:          "delete",
		Short:        "Delete a workspaces",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := deleteWorkspacesOptions.Complete(args); err != nil {
				return err
			}

			if err := deleteWorkspacesOptions.Validate(); err != nil {
				return err
			}

			return deleteWorkspacesOptions.Run(c.Context())
		},
	}

	useWorkspacesOptions := plugin.NewUseWorkspacesOptions(streams)
	useWorkspacesCmd := &cobra.Command{
		Use:          "use",
		Short:        "Use a workspaces",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := useWorkspacesOptions.Complete(args); err != nil {
				return err
			}

			if err := useWorkspacesOptions.Validate(); err != nil {
				return err
			}

			return useWorkspacesOptions.Run(c.Context())
		},
	}

	getWorkspacesOptions.BindFlags(getWorkspacesCmd)
	cmd.AddCommand(getWorkspacesCmd)

	createWorkspacesOptions.BindFlags(createWorkspacesCmd)
	cmd.AddCommand(createWorkspacesCmd)

	deleteWorkspacesOptions.BindFlags(deleteWorkspacesCmd)
	cmd.AddCommand(deleteWorkspacesCmd)

	useWorkspacesOptions.BindFlags(useWorkspacesCmd)
	cmd.AddCommand(useWorkspacesCmd)

	return cmd, nil
}
