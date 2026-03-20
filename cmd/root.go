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
		// If user explicitly asks for help, use Cobra's help system
		helpFlag, _ := cmd.Flags().GetBool("help")
		if helpFlag {
			cmd.Help()
			return
		}

		printer.PrintBanner()

		printer.Section(messages.RootAvailableCommands)
		printer.Command(messages.RootCmdDoctor)
		printer.Command(messages.RootCmdInstall)
		printer.Command(messages.RootCmdList)
		printer.Command(messages.RootCmdGraph)
		printer.Command(messages.RootCmdVersion)

		printer.Section(messages.RootTypicalWorkflow)
		printer.Command(messages.RootWorkflowDoctor)
		printer.Command(messages.RootWorkflowInstall)
		printer.Command(messages.RootWorkflowInstallAll)
		printer.Blank()
	},
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
