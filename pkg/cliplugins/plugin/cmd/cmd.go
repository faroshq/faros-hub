package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/faroshq/faros-hub/pkg/cliplugins/plugin/plugin"
)

// New provides a cobra command for workload operations.
func New(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Aliases:          []string{"plugins"},
		Use:              "plugin",
		Short:            "Manages plugins",
		SilenceUsage:     true,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	getOptions := plugin.NewGetOptions(streams)
	getCmd := &cobra.Command{
		Use:          "get",
		Short:        "Get a plugins",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := getOptions.Complete(args); err != nil {
				return err
			}

			if err := getOptions.Validate(); err != nil {
				return err
			}

			return getOptions.Run(c.Context())
		},
	}

	enableOptions := plugin.NewEnableOptions(streams)
	enableCmd := &cobra.Command{
		Use:          "enable",
		Short:        "Enable a plugins",
		Aliases:      []string{"request"},
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := enableOptions.Complete(args); err != nil {
				return err
			}

			if err := enableOptions.Validate(); err != nil {
				return err
			}

			return enableOptions.Run(c.Context())
		},
	}

	enableOptions.BindFlags(enableCmd)
	cmd.AddCommand(enableCmd)

	getOptions.BindFlags(getCmd)
	cmd.AddCommand(getCmd)

	return cmd, nil
}
