// Based on https://github.com/kcp-dev/kcp/

package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"
	cliflag "k8s.io/component-base/cli/flag"
)

const (
	usageFmt = "Usage:\n  %s\n"
)

// setPartialUsageAndHelpFunc set both usage and help function.
// Print the flag sets we need instead of all of them.
func setPartialUsageAndHelpFunc(cmd *cobra.Command, fss cliflag.NamedFlagSets, cols int, flags []string) {
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		printMostImportantFlags(cmd.OutOrStderr(), fss, cols, flags)
		fmt.Fprintf(cmd.OutOrStderr(), "\nUse \"%s\" for a list of all flags available.\n", cmd.CommandPath())
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		printMostImportantFlags(cmd.OutOrStdout(), fss, cols, flags)
		fmt.Fprintf(cmd.OutOrStderr(), "\nUse \"%s options\" for a list of all flags available.\n", cmd.CommandPath())
	})
}

func printMostImportantFlags(w io.Writer, fss cliflag.NamedFlagSets, cols int, visibleFlags []string) {
	visibleFlagsSet := sets.NewString(visibleFlags...)
	filteredFFS := cliflag.NamedFlagSets{}
	filteredFS := filteredFFS.FlagSet("Most important")

	for _, name := range fss.Order {
		fs := fss.FlagSets[name]
		if !fs.HasFlags() {
			continue
		}

		fs.VisitAll(func(f *pflag.Flag) {
			if visibleFlagsSet.Has(f.Name) {
				filteredFS.AddFlag(f)
			}
		})
	}

	cliflag.PrintSections(w, filteredFFS, cols)
}
