package cmd

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Args:  cobra.NoArgs,
	Short: messages.CmdSetupShortDescription,
	Long:  messages.CmdSetupLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		printer.Section(messages.EnvironmentSetup)

		if _, err := environment.SetupInteractive(); err != nil {
			printer.Fail(messages.EnvironmentSetupFailed)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
