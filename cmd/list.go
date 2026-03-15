package cmd

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/tools"
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

The output shows the available tools along with a brief description
so you can decide which tools to install individually or as part
of the full development stack.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Available tools:")

		for _, t := range tools.List() {
			fmt.Printf("%s - %s\n", t.Name(), t.Description())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
