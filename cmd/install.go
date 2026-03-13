package cmd

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [tool]",
	Args:  cobra.ExactArgs(1),
	Short: "Install a tool",
	Long: `
	Install a development tool into your local environment.

This command installs a supported tool and prepares it for use with the
current mdev environment configuration. The tool will be downloaded,
installed, and configured using the paths and settings defined in your
mdev environment.

Before running this command, your environment must be initialized and
validated using 'mdev doctor'. The install process depends on the
configured directories, tool paths, and system checks performed during
that step.

If the tool is already installed, the command will detect it and skip
the installation to avoid overwriting an existing setup.

Usage:
  mdev install [tool]

Arguments:
  tool    Name of the tool to install.

Behavior:
  • Validates that the environment is configured.
  • Checks whether the requested tool is supported.
  • Detects if the tool is already installed.
  • Runs the tool-specific installation process.

Examples:
  mdev install java
  mdev install gradle
  mdev install maven

Notes:
  Each tool provides its own installation logic. The command acts as a
  dispatcher that resolves the requested tool and executes its install
  routine using the current environment configuration.
	`,
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

		if tool.IsInstalled(env) {
			fmt.Println("✓", name, "already installed")
			return
		}

		fmt.Println("Installing", name)

		// install
		err = tool.Install(env)
		if err != nil {
			fmt.Println("Installation failed:", err)
			return
		}

		err = tool.Configure(env)
		if err != nil {
			fmt.Println("Configuration failed:", err)
			return
		}

		err = tool.Verify(env)
		if err != nil {
			fmt.Println("Verification failed:", err)
			return
		}

		fmt.Println("✓ Installed and configured", name)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
