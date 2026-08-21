package cmd

import (
	"github.com/jamesawo/mdev/internal/command/install"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

// installCmd represents the installation command
var installCmd = &cobra.Command{
	Use:   "install [tool]",
	Args:  cobra.MaximumNArgs(1),
	Short: messages.InstallCmdShortDescription,
	Long:  messages.InstallCmdLongDescription,
	Run: func(cmd *cobra.Command, args []string) {

		// Prevent invalid usage
		if installAll && len(args) > 0 {
			printer.Fail(messages.InstallAllWithToolError)
			return
		}

		env, err := environment.FromConfig()

		if err != nil {
			printer.Fail("mdev is not configured")
			printer.Info(messages.CommonRun + " " + messages.SetupCommand)
			return
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
			name = name + messages.ToolsInstalledSuffix
		}

		options = append(options, name)
		toolMap[name] = t.Name()
	}

	selected, err := interactive.MultiSelect(
		printer.FormatIndent(1, messages.InstallSelectTools),
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
