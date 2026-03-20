package services

import exec "github.com/jamesawo/mdev/internal/infrastructure/exec"

// Default returns all services used by mdev
func Default() []Service {

	// Create one runner instance
	runner := &exec.CommandRunner{}

	return []Service{
		NewOllama(runner),
		NewPodman(runner),
	}
}
