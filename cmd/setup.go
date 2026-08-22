package cmd

import (
	"errors"

	commandsetup "github.com/jamesawo/mdev/internal/command/setup"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var setupStoragePath string

var defaultRunSetup = commandsetup.Run
var runSetup = defaultRunSetup

var setupCmd = &cobra.Command{
	Use:   messages.SetupCommandName,
	Args:  cobra.NoArgs,
	Short: messages.SetupCmdShortDescription,
	Long:  messages.SetupCmdLongDescription,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed(messages.SetupYesFlagName) || cmd.InheritedFlags().Changed(messages.SetupYesFlagName) {
			return errors.New(messages.SetupYesUnsupported)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetup(commandsetup.Options{StoragePath: setupStoragePath})
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVar(&setupStoragePath, messages.SetupStoragePathFlagName, "", messages.SetupStoragePathFlag)
	defaultHelp := setupCmd.HelpFunc()
	setupCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		flag := cmd.InheritedFlags().Lookup(messages.SetupYesFlagName)
		if flag == nil {
			defaultHelp(cmd, args)
			return
		}
		hidden := flag.Hidden
		flag.Hidden = true
		defer func() { flag.Hidden = hidden }()
		defaultHelp(cmd, args)
	})
}
