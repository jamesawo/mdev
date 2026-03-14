package cmd

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [tool]",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your command",
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
			fmt.Println("Environment not configured. Run `mdev doctor` first.")
			return
		}

		tool, ok := tools.Get(name)
		if !ok {
			fmt.Println("Unknown tool:", name)
			return
		}

		if !confirmUninstall(name) {
			fmt.Println("Cancelled.")
			return
		}

		fmt.Println("Uninstalling", name)

		err = tool.Uninstall(env)
		if err != nil {
			fmt.Println("Uninstall failed:", err)
			return
		}

		fmt.Println("✓ Uninstalled", name)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func confirmUninstall(name string) bool {

	fmt.Printf("Are you sure you want to uninstall %s? (y/N): ", name)

	var input string
	fmt.Scanln(&input)

	return input == "y" || input == "Y"
}
