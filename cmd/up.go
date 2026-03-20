package cmd

import (
	"context"

	exec "github.com/jamesawo/mdev/internal/infrastructure/exec"
	"github.com/jamesawo/mdev/internal/services"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start development environment",
	Long: `Start the development environment.

mdev up brings your development environment by starting all required runtime services.

Examples:
  mdev up

This typically includes:
- starting podman machine
- starting ollama service

Use this at the start of your development session.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		runner := &exec.CommandRunner{}
		svcs := services.Default()

		manager := services.New(svcs, runner)

		return manager.Up(context.Background())
	},
}
