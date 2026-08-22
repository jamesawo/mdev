package cmd

import (
	"context"
	"os"
	"os/signal"

	// Register tools
	_ "github.com/jamesawo/mdev/internal/tools/gradle"
	_ "github.com/jamesawo/mdev/internal/tools/java"
	_ "github.com/jamesawo/mdev/internal/tools/maven"
	_ "github.com/jamesawo/mdev/internal/tools/nvm"
	_ "github.com/jamesawo/mdev/internal/tools/podman"
	_ "github.com/jamesawo/mdev/internal/tools/sdkman"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/messages"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "mdev",
	Short:        messages.RootCmdShortDescription,
	Long:         messages.RootCmdLongDescription,
	Version:      versionMetadata.Version,
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		confirmation.Configure(confirmAll)
	},
}

var confirmAll bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&confirmAll, "yes", "y", false, messages.RootFlagConfirmAll)
	rootCmd.SetVersionTemplate(messages.RootVersionTemplate)
}

// Execute runs the CLI.
func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
