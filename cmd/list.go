package cmd

import (
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported development tools",
	Long: `List all development tools supported by mdev.

This command displays the tools that mdev knows how to manage.
Each tool includes a short description and can be installed,
configured, and managed through the mdev lifecycle.

Typical usage:

  mdev list
`,
	Run: func(cmd *cobra.Command, args []string) {

		env, _ := environment.FromConfig()

		printer.Section("Available tools")

		for _, t := range tools.List() {

			name := t.Name()

			if env != nil && t.IsInstalled(env) {
				printer.Success(name + " (installed)")
				continue
			}

			printer.Fail(name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
