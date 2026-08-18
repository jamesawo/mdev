package cmd

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

const (
	version = "0.1.0"
	author  = "James Aworo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: messages.VersionCmdShortDescription,
	Long:  messages.VersionCmdLongDescription,

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(messages.VersionInfo(version))
		fmt.Println(messages.VersionAuthor(author))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
