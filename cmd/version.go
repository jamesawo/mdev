package cmd

import (
	"io"

	commandversion "github.com/jamesawo/mdev/internal/command/version"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var versionMetadata = commandversion.DefaultMetadata()

var defaultRunVersion = func(out io.Writer, metadata commandversion.Metadata, options commandversion.Options) error {
	return commandversion.Run(out, metadata, options)
}

var runVersion = defaultRunVersion

var versionCmd = &cobra.Command{
	Use:   messages.VersionCommandName,
	Args:  cobra.NoArgs,
	Short: messages.VersionCmdShortDescription,
	Long:  messages.VersionCmdLongDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVersion(cmd.OutOrStdout(), versionMetadata, commandversion.Options{JSON: versionJSON})
	},
}

var versionJSON bool

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&versionJSON, messages.VersionJSONFlagName, false, messages.VersionJSONFlag)
}
