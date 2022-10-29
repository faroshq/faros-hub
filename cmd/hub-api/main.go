// Based on https://github.com/kcp-dev/kcp/

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/cmd/help"
	"github.com/kcp-dev/kcp/pkg/embeddedetcd"
	kcpfeatures "github.com/kcp-dev/kcp/pkg/features"
	"github.com/kcp-dev/kcp/pkg/server"
	"github.com/kcp-dev/kcp/pkg/server/options"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	genericapiserver "k8s.io/apiserver/pkg/server"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/rest"
	"k8s.io/component-base/cli"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/config"
	"k8s.io/component-base/logs"
	"k8s.io/component-base/term"
	"k8s.io/component-base/version"
	"k8s.io/klog"

	"github.com/faroshq/faros-hub/pkg/bootstrap/templates/root"
	farosconfig "github.com/faroshq/faros-hub/pkg/config"
	"github.com/faroshq/faros-hub/pkg/controllers"
	faroserver "github.com/faroshq/faros-hub/pkg/server"
	bootstraputils "github.com/faroshq/faros-hub/pkg/util/bootstrap"
)

var allInOne bool

func main() {
	cmd := &cobra.Command{
		Use:   "hub",
		Short: "Faros hub",
		Long: help.Doc(`
			Faros hub will start minimal faros hub.
		`),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())

	// manually extract root directory from flags first as it influence all other flags
	rootDir := ".faros"
	for i, f := range os.Args {
		if f == "--root-directory" {
			if i < len(os.Args)-1 {
				rootDir = os.Args[i+1]
			} // else let normal flag processing fail
		} else if strings.HasPrefix(f, "--root-directory=") {
			rootDir = strings.TrimPrefix(f, "--root-directory=")
		}
	}

	serverOptions := options.NewOptions(rootDir)
	serverOptions.GenericControlPlane.Logs.Config.Verbosity = config.VerbosityLevel(12)
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the control plane process",
		Long: help.Doc(`
			Start the control plane process
		`),
		PersistentPreRunE: func(*cobra.Command, []string) error {
			// silence client-go warnings.
			// apiserver loopback clients should not log self-issued warnings.
			rest.SetDefaultWarningHandler(rest.NoWarnings{})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// run as early as possible to avoid races later when some components (e.g. grpc) start early using klog
			if err := serverOptions.GenericControlPlane.Logs.ValidateAndApply(kcpfeatures.DefaultFeatureGate); err != nil {
				return err
			}

			c, err := farosconfig.LoadController()
			if err != nil {
				return err
			}

			t, err := faroserver.New(c)
			if err != nil {
				return err
			}

			serverOptions.Extra.AdditionalAPIHandlers = t.GetHandlers()

			completed, err := serverOptions.Complete()
			if err != nil {
				return err
			}

			//completed.Extra.BatteriesIncluded = []string{}
			klog.Infof("Batteries included: %s", strings.Join(completed.Extra.BatteriesIncluded, ","))

			if errs := completed.Validate(); len(errs) > 0 {
				return errors.NewAggregate(errs)
			}

			config, err := server.NewConfig(completed)
			if err != nil {
				return err
			}

			// set tunnels true
			runtime.Must(utilfeature.DefaultMutableFeatureGate.Set(fmt.Sprintf("%s=true", kcpfeatures.SyncerTunnel)))

			completedConfig, err := config.Complete()
			if err != nil {
				return err
			}

			ctx := genericapiserver.SetupSignalContext()

			// the etcd server must be up before NewServer because storage decorators access it right away
			if completedConfig.EmbeddedEtcd.Config != nil {
				if err := embeddedetcd.NewServer(completedConfig.EmbeddedEtcd).Run(ctx); err != nil {
					return err
				}
			}

			s, err := server.NewServer(completedConfig)
			if err != nil {
				return err
			}

			// Add hook to populate tunnels clients
			// Register a post-start hook that connects to the api-server
			s.AddPostStartHook("connect-to-api", func(ctx genericapiserver.PostStartHookContext) error {
				// Create a new client using the client config from our newly created api-server
				err := t.SeedClients(ctx.LoopbackClientConfig)
				return err
			})

			// based on https://github.com/kcp-dev/kcp/blob/124785e6e911c0d564db08667b985718a37c18f2/pkg/server/server.go
			// Add hook for root bootstrap
			s.AddPostStartHook("bootstrap-faros-workspace-types", func(ctx genericapiserver.PostStartHookContext) error {
				if s.Options.Extra.ShardName == tenancyv1alpha1.RootShard {
					if err := root.Bootstrap(goContext(ctx),
						s.BootstrapApiExtensionsClusterClient.Cluster(tenancyv1alpha1.RootCluster).Discovery(),
						s.BootstrapDynamicClusterClient.Cluster(tenancyv1alpha1.RootCluster),
						bootstraputils.ReplaceOption()); err != nil {
						return err // don't klog.Fatal. This only happens when context is cancelled.
					}
				}
				return nil
			})

			// Start the controllers in a goroutine
			// Would be better not to do this in production
			if allInOne {
				// Add hook to start controllers too
				s.AddPostStartHook("run-faros-controllers", func(ctxP genericapiserver.PostStartHookContext) error {

					c.RestConfig = ctxP.LoopbackClientConfig

					controllers, err := controllers.New(c)
					if err != nil {
						return err
					}

					go func() {
						for {
							err := controllers.Run(ctx)
							if err != nil {
								klog.Error(err)
							}
							time.Sleep(5 * time.Second)
						}
					}()
					return nil
				})
			}

			return s.Run(ctx)
		},
	}

	// add start named flag sets to start flags
	namedStartFlagSets := serverOptions.Flags()

	globalflag.AddGlobalFlags(namedStartFlagSets.FlagSet("global"), cmd.Name(), logs.SkipLoggingConfigurationFlags())
	startFlags := startCmd.Flags()
	startFlags.BoolVarP(&allInOne, "all-in-one", "a", false, "Should start all-in-one")

	for _, f := range namedStartFlagSets.FlagSets {
		startFlags.AddFlagSet(f)
	}

	startOptionsCmd := &cobra.Command{
		Use:   "options",
		Short: "Show all start command options",
		Long: help.Doc(`
			Show all start command options

			"faros start"" has a large number of options. This command shows all of them.
		`),
		PersistentPreRunE: func(*cobra.Command, []string) error {
			// silence client-go warnings.
			// apiserver loopback clients should not log self-issued warnings.
			rest.SetDefaultWarningHandler(rest.NoWarnings{})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), usageFmt, startCmd.UseLine())
			cliflag.PrintSections(cmd.OutOrStderr(), namedStartFlagSets, cols)
			return nil
		},
	}
	startCmd.AddCommand(startOptionsCmd)
	cmd.AddCommand(startCmd)

	setPartialUsageAndHelpFunc(startCmd, namedStartFlagSets, cols, []string{
		"etcd-servers",
		"batteries-included",
		"run-virtual-workspaces",
	})

	help.FitTerminal(cmd.OutOrStdout())

	if v := version.Get().String(); len(v) == 0 {
		cmd.Version = "<unknown>"
	} else {
		cmd.Version = v
	}
	os.Exit(cli.Run(cmd))
}

// goContext turns the PostStartHookContext into a context.Context for use in routines that may or may not
// run inside of a post-start-hook. The k8s APIServer wrote the post-start-hook context code before contexts
// were part of the Go stdlib.
func goContext(parent genericapiserver.PostStartHookContext) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func(done <-chan struct{}) {
		<-done
		cancel()
	}(parent.StopCh)
	return ctx
}
