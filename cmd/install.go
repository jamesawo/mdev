package cmd

import (
	"errors"

	commandinstall "github.com/jamesawo/mdev/internal/command/install"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var defaultRunInstall = commandinstall.Run
var runInstall = defaultRunInstall

var installCmd = &cobra.Command{
	Use:   messages.InstallCommandUse,
	Args:  cobra.MaximumNArgs(1),
	Short: messages.InstallCmdShortDescription,
	Long:  messages.InstallCmdLongDescription,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if installAll && len(args) > 0 {
			return errors.New(messages.InstallAllWithToolError)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tool := ""
		if len(args) == 1 {
			tool = args[0]
		}
		return runInstall(cmd.Context(), commandinstall.Streams{
			In: cmd.InOrStdin(), Out: cmd.OutOrStdout(), Err: cmd.ErrOrStderr(),
		}, commandinstall.Options{Tool: tool, All: installAll, AssumeYes: confirmAll})
	},
}

var installAll bool

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&installAll, messages.InstallAllFlagName, false, messages.InstallAllFlag)
}
