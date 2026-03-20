package services

import (
	"context"
	"fmt"

	exec "github.com/jamesawo/mdev/internal/infrastructure/exec"
)

// Manager orchestrates all services
type Manager struct {
	services []Service
	runner   exec.Runner
}

// New Constructor
func New(services []Service, runner exec.Runner) *Manager {
	return &Manager{
		services: services,
		runner:   runner,
	}
}

// Up - starts all services
func (m *Manager) Up(ctx context.Context) error {

	for _, s := range m.services {
		if err := s.Start(ctx); err != nil {
			return fmt.Errorf("%s failed to start: %w", s.Name(), err)
		}
	}

	return nil
}

// Down stops services and unmounts disk
func (m *Manager) Down(ctx context.Context) error {

	for _, s := range m.services {
		if err := s.Stop(ctx); err != nil {
			return fmt.Errorf("%s failed to stop: %w", s.Name(), err)
		}
	}

	// Unmount SSD
	return m.runner.Run(ctx, "diskutil", "unmount", "/Volumes/scandisk")
}
