package cmd

import (
	"io"

	"github.com/jamesawo/mdev/internal/command/list"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var defaultRunList = func(out io.Writer) error {
	return list.Run(out)
}

var runList = defaultRunList

var listCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.NoArgs,
	Short: messages.ListCmdShortDescription,
	Long:  messages.ListCmdLongDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
