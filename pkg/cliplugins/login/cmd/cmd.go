package cmd

import (
	"github.com/faroshq/faros-hub/pkg/cliplugins/login/plugin"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// New provides a cobra command for login
func New(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	loginOptions := plugin.NewLoginSetupOptions(streams)

	cmd := &cobra.Command{
		Aliases:          []string{"login"},
		Use:              "login",
		Short:            "Manages Faros access via SSO",
		SilenceUsage:     true,
		TraverseChildren: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := loginOptions.Complete(args); err != nil {
				return err
			}

			if err := loginOptions.Validate(); err != nil {
				return err
			}

			return loginOptions.Run(c.Context())
		},
	}
	loginOptions.BindFlags(cmd)

	return cmd, nil
}
