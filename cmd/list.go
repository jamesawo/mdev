package cmd

import (
	"github.com/jamesawo/mdev/internal/command/list"
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
		list.Run()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
