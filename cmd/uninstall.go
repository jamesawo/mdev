package cmd

import (
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/jamesawo/mdev/internal/uninstall"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [tool]",
	Args:  cobra.ExactArgs(1),
	Short: "Uninstall a tool from your local environment",
	Long: `Remove an installed tool from the local environment.

This command uninstalls a tool that was previously installed using
mdev. It removes the tool binaries and any managed directories that
belong to the mdev environment while keeping unrelated user files
untouched.

The command validates that the environment is configured and that
the specified tool is known by mdev before attempting removal.

Usage:
  mdev uninstall [tool]

Arguments:
  tool    Name of the tool to uninstall.

Behavior:
  • Verifies the mdev environment configuration.
  • Confirms the tool is supported by mdev.
  • Removes the installed tool from the managed environment.

Examples:
  mdev uninstall java
  mdev uninstall gradle
  mdev uninstall maven

Notes:
  Only tools managed by mdev can be removed using this command.
  If the tool is not installed, the command will exit without
  making changes.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		env, err := environment.FromConfig()
		if err != nil {
			printer.Fail("Environment not configured. Run `mdev doctor` first.")
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
