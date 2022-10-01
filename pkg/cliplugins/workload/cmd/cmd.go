// Based on https://github.com/kcp-dev/kcp/tree/main/pkg/cliplugins/workload

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/faroshq/faros-hub/pkg/cliplugins/workload/plugin"
)

var (
	syncExample = `
	# Ensure a syncer is running on the specified sync target.
	%[1]s workload sync <sync-target-name> --syncer-image <kcp-syncer-image> -o syncer.yaml
	KUBECONFIG=<pcluster-config> kubectl apply -f syncer.yaml

	# Directly apply the manifest
	%[1]s workload sync <sync-target-name> --syncer-image <kcp-syncer-image> -o - | KUBECONFIG=<pcluster-config> kubectl apply -f -`
)

// New provides a cobra command for workload operations.
func New(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Aliases:          []string{"workloads"},
		Use:              "workload",
		Short:            "Manages cluster sync targets",
		SilenceUsage:     true,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Sync command
	syncOptions := plugin.NewSyncOptions(streams)

	enableSyncerCmd := &cobra.Command{
		Use:          "sync <sync-target-name> --syncer-image <kcp-syncer-image> [--resources=<resource1>,<resource2>..] -o <output-file>",
		Short:        "Create a synctarget in kcp with service account and RBAC permissions. Output a manifest to deploy a syncer for the given sync target in a physical cluster.",
		Example:      fmt.Sprintf(syncExample, "kubectl kcp"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return c.Help()
			}

			if err := syncOptions.Complete(args); err != nil {
				return err
			}

			if err := syncOptions.Validate(); err != nil {
				return err
			}

			return syncOptions.Run(c.Context())
		},
	}

	syncOptions.BindFlags(enableSyncerCmd)
	cmd.AddCommand(enableSyncerCmd)

	return cmd, nil
}
