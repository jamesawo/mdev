package cmd

import (
	"context"

	upcmd "github.com/jamesawo/mdev/internal/command/up"
	"github.com/jamesawo/mdev/internal/ui/messages"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: messages.CmdUpShortDescription,
	Long:  messages.CmdUpLongDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return upcmd.New().Execute(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
