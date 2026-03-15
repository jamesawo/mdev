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

	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdev",
	Short: "Automate development environment setup on macOS",
	Long: `mdev is a command-line tool for setting up and managing
a development environment on macOS.

It installs development tools, configures them, and relocates large
tool caches to external storage to keep your system disk clean.`,

	Run: func(cmd *cobra.Command, args []string) {
		printer.PrintBanner()

		printer.Section("Available commands")
		printer.Command("mdev doctor   Initialize and validate your environment")
		printer.Command("mdev install  Install development tools")
		printer.Command("mdev list     Show supported tools and their status")
		printer.Command("mdev graph    Show dependency graph between tools")

		printer.Section("Typical workflow")
		printer.Command("mdev doctor")
		printer.Command("mdev install")
		printer.Command("mdev install --all")
	},
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
