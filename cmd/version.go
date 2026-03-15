package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "0.1.0"
	author  = "James Aworo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show mdev version information",
	Long: `Display the current version of mdev and basic
information about the project.`,

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("mdev", version)
		fmt.Println("Created by", author)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
