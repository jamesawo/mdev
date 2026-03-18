package exec

import (
	"context"
	"os"
	"os/exec"
)

// Runner defines how commands are executed.
// We use an interface so we can later mock it in tests.
type Runner interface {
	Run(ctx context.Context, name string, args ...string) error
}

// CommandRunner is the real implementation that runs commands on the system.
type CommandRunner struct{}

// Run executes a system command like:
// brew install go
// podman machine start
func (r *CommandRunner) Run(ctx context.Context, name string, args ...string) error {

	// Create command with context (allows future timeout/cancel)
	cmd := exec.CommandContext(ctx, name, args...)

	// Attach stdout/stderr so user sees output in terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}
