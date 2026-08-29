package command

import (
	"context"

	"github.com/jamesawo/mdev/internal/infrastructure/subprocess"
)

// Run executes a subprocess with the process standard streams.
func Run(name string, args ...string) error {
	return RunContext(context.Background(), name, args...)
}

// RunContext executes a subprocess with standard streams and cancellation.
func RunContext(ctx context.Context, name string, args ...string) error {
	return subprocess.Run(ctx, name, args...)
}
