package up

import (
	"context"

	exec "github.com/jamesawo/mdev/internal/infrastructure/exec"
	"github.com/jamesawo/mdev/internal/services"
)

type Command struct{}

func New() *Command {
	return &Command{}
}

func (c *Command) Execute(ctx context.Context) error {
	runner := &exec.CommandRunner{}
	svcs := services.Default()

	manager := services.New(svcs, runner)

	return manager.Up(ctx)
}
