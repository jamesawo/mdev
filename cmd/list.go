package cmd

import (
	"io"

	"github.com/jamesawo/mdev/internal/command/list"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var defaultRunList = func(out io.Writer, options list.Options) error {
	return list.Run(out, options)
}

var runList = defaultRunList

var listCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.NoArgs,
	Short: messages.ListCmdShortDescription,
	Long:  messages.ListCmdLongDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList(cmd.OutOrStdout(), list.Options{JSON: listJSON})
	},
}

var listJSON bool

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&listJSON, "json", false, messages.ListJSONFlag)
}
