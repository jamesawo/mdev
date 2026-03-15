package cmd

import (
	"github.com/jamesawo/mdev/internal/command/install"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [tool]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Install a tool in your local environment.",
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

		// Prevent invalid usage
		if installAll && len(args) > 0 {
			printer.Fail("cannot use --all with a specific tool")
			return
		}

		env, err := environment.FromConfig()

		// First-time setup
		if err != nil {

			printer.Section("Environment setup")
			printer.Info("Choose where mdev should store development tool data.")

			if !interactive.AskYesNo("Create the directory now?") {
				printer.Info("Aborted.")
				printer.Blank()
				return
			}

			env, err = environment.SetupInteractive()
			if err != nil {
				printer.Fail("environment setup failed")
				printer.Blank()
				return
			}

			printer.Success("Environment initialized")
		}

		// install all tools
		if installAll {

			err := install.RunAll(env)
			if err != nil {
				printer.Fail(err.Error())
			}

			return
		}

		// install a single tool
		if len(args) == 1 {

			err := install.RunSingle(env, args[0])
			if err != nil {
				printer.Fail(err.Error())
			}

			return
		}

		// interactive mode
		runInteractiveInstall(env)

	},
}
var installAll bool

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&installAll, "all", false, "Install all tools")
}

func runInteractiveInstall(env *environment.Environment) {

	var options []string
	toolMap := map[string]string{}

	for _, t := range tools.List() {

		name := t.Name()

		if t.IsInstalled(env) {
			name = name + " (installed)"
		}

		options = append(options, name)
		toolMap[name] = t.Name()
	}

	selected, err := interactive.MultiSelect(
		"Select tools to install",
		options,
	)

	if err != nil {
		printer.Fail(err.Error())
		return
	}

	var names []string

	for _, s := range selected {
		names = append(names, toolMap[s])
	}

	err = install.RunSelection(env, names)
	if err != nil {
		printer.Fail(err.Error())
	}
}
