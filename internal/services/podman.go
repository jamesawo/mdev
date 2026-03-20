package services

import (
	"context"

	exec "github.com/jamesawo/mdev/internal/infrastructure/exec"
)

// Podman represents the Podman machine service
type Podman struct {
	runner exec.Runner
}

// NewPodman Constructor
func NewPodman(r exec.Runner) Podman {
	return Podman{runner: r}
}

func (p Podman) Name() string {
	return "podman"
}

// Start runs: podman machine start
func (p Podman) Start(ctx context.Context) error {
	return p.runner.Run(ctx, "podman", "machine", "start")
}

// Stop runs: podman machine stop
func (p Podman) Stop(ctx context.Context) error {
	return p.runner.Run(ctx, "podman", "machine", "stop")
}

func (p Podman) Status(ctx context.Context) Status {
	return Status{Running: false} // implement later
}
