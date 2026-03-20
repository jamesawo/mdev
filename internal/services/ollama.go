package services

import (
	"context"

	exec "github.com/jamesawo/mdev/internal/infrastructure/exec"
)

// Ollama represents the Ollama service (local LLM server)
type Ollama struct {
	runner exec.Runner
}

// NewOllama Constructor
func NewOllama(r exec.Runner) Ollama {
	return Ollama{runner: r}
}

func (o Ollama) Name() string {
	return "ollama"
}

// Start runs: brew services start ollama
func (o Ollama) Start(ctx context.Context) error {
	return o.runner.Run(ctx, "brew", "services", "start", "ollama")
}

// Stop runs: brew services stop ollama
func (o Ollama) Stop(ctx context.Context) error {
	return o.runner.Run(ctx, "brew", "services", "stop", "ollama")
}

func (o Ollama) Status(ctx context.Context) Status {
	return Status{Running: false} // implement later
}
