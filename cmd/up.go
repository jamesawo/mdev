package cmd

import (
	"context"

	upcmd "github.com/jamesawo/mdev/internal/command/up"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start development services (like podman, ollama)",
	Long: `Start the development environment.

mdev up starts the runtime services and tools required during development.

This typically includes:
- podman machine
- ollama service


Use this command at the beginning of your development session.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return upcmd.New().Execute(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
