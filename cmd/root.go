package cmd

import (
	"os"

	// Register tools
	_ "github.com/jamesawo/mdev/internal/tools/gradle"
	_ "github.com/jamesawo/mdev/internal/tools/java"
	_ "github.com/jamesawo/mdev/internal/tools/maven"
	_ "github.com/jamesawo/mdev/internal/tools/nvm"
	_ "github.com/jamesawo/mdev/internal/tools/podman"
	_ "github.com/jamesawo/mdev/internal/tools/sdkman"
	"github.com/jamesawo/mdev/internal/ui/messages"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdev",
	Short: messages.CmdRootShortDescription,
	Long:  messages.CmdRootLongDescription,
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
