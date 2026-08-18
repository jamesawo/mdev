package cmd

import (
	"github.com/jamesawo/mdev/internal/command/list"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: messages.ListCmdShortDescription,
	Long:  messages.ListCmdLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		list.Run()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
