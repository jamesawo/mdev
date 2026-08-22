package command

import (
	"context"
	"os"
	"os/exec"
)

// Run executes a subprocess with the process standard streams.
func Run(name string, args ...string) error {
	return RunContext(context.Background(), name, args...)
}

// RunContext executes a subprocess with standard streams and cancellation.
func RunContext(ctx context.Context, name string, args ...string) error {

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
