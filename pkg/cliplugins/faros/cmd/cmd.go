/*
Copyright 2022 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	workloadcmd "github.com/faroshq/faros-hub/pkg/cliplugins/workload/cmd"
	workspacecmd "github.com/kcp-dev/kcp/pkg/cliplugins/workspace/cmd"
)

// New returns a cobra.Command for faros actions.
func New(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	workspaceCmd, err := workspacecmd.New(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	workloadCmd, err := workloadcmd.New(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "faros",
		Short: "Manage faros",
	}

	cmd.AddCommand(workspaceCmd)
	cmd.AddCommand(workloadCmd)

	return cmd, nil
}
