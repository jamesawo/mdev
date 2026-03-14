package system

import (
	"github.com/jamesawo/mdev/internal/command"
)

type Brew struct{}

func (b *Brew) Name() string {
	return "brew"
}

func (b *Brew) Check() bool {
	return CommandExists("brew")
}

func (b *Brew) Install() error {

	return command.Run(
		"bash",
		"-c",
		`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
	)
}

func init() {
	Register(&Brew{})
}
