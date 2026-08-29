package packagemanager

import (
	"context"

	"github.com/jamesawo/mdev/internal/infrastructure/subprocess"
)

// CommandRunner executes package-manager commands with process output streams.
type CommandRunner struct{}

// Run executes a package-manager command without caller cancellation.
func (r *CommandRunner) Run(name string, args ...string) error {
	return r.RunContext(context.Background(), name, args...)
}

// RunContext executes a package-manager command with caller cancellation.
func (r *CommandRunner) RunContext(ctx context.Context, name string, args ...string) error {
	return subprocess.Run(ctx, name, args...)
}
