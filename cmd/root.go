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

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdev",
	Short: "Automate development environment setup on macOS",
	Long: `mdev is a command-line tool for setting up and managing
a development environment on macOS.

It installs development tools, configures them, and relocates large
tool caches to external storage to keep your system disk clean.`,
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
