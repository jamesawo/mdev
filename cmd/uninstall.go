package cmd

import (
	"github.com/jamesawo/mdev/internal/command/uninstall"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [tool]",
	Args:  cobra.ExactArgs(1),
	Short: messages.CmdUninstallShortDescription,
	Long:  messages.CmdUninstallLongDescription,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		env, err := environment.FromConfig()
		if err != nil {
			printer.Fail(messages.ErrEnvironmentNotConfigured)
			return
		}

		err = uninstall.Run(env, name)
		if err != nil {
			printer.Fail(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
