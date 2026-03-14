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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
