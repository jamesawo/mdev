package exec

import (
	"context"

	"github.com/jamesawo/mdev/internal/infrastructure/subprocess"
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
	return subprocess.Run(ctx, name, args...)
}
