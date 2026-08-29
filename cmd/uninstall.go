package cmd

import (
	commanduninstall "github.com/jamesawo/mdev/internal/command/uninstall"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var defaultRunUninstall = commanduninstall.Run
var runUninstall = defaultRunUninstall

var uninstallCmd = &cobra.Command{
	Use:   messages.UninstallCommandUse,
	Args:  cobra.ExactArgs(1),
	Short: messages.UninstallCmdShortDescription,
	Long:  messages.UninstallCmdLongDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUninstall(cmd.Context(), commanduninstall.Streams{
			In:  cmd.InOrStdin(),
			Out: cmd.OutOrStdout(),
			Err: cmd.ErrOrStderr(),
		}, commanduninstall.Options{Tool: args[0], AssumeYes: confirmAll})
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
