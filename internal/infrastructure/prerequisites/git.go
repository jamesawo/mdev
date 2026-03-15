package prerequisites

import "github.com/jamesawo/mdev/internal/command"

type Git struct{}

func (g *Git) Name() string {
	return "git"
}

func (g *Git) Check() bool {
	return CommandExists("git")
}

func (g *Git) Install() error {
	return command.Run("brew", "install", "git")
}

func init() {
	Register(&Git{})
}
