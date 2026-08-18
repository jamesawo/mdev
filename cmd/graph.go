package cmd

import (
	"github.com/jamesawo/mdev/internal/command/graph"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: messages.GraphCmdShortDescription,
	Long:  messages.GraphCmdLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		graph.Run()
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
