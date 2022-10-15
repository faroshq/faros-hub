package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	agentcmd "github.com/faroshq/faros-hub/pkg/cliplugins/agent/cmd"
)

// New returns a cobra.Command for faros actions.
func New(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	agentCmd, err := agentcmd.New(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "faros",
		Short: "Manage faros",
	}

	cmd.AddCommand(agentCmd)

	return cmd, nil
}
