package cmd

import (
	"errors"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Args:  cobra.NoArgs,
	Short: messages.SetupCmdShortDescription,
	Long:  messages.SetupCmdLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		printer.Section(messages.SetupTitle)

		if _, err := environment.SetupInteractive(); errors.Is(err, environment.ErrSetupCancelled) {
			printer.Info(messages.CommonAborted)
		} else if err != nil {
			printer.Fail(messages.SetupFailed)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
