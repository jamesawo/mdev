package system

import (
	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/fs"
)

type Brew struct{}

func (b *Brew) Name() string {
	return "brew"
}

func (b *Brew) Check() bool {

	if CommandExists("brew") {
		return true
	}

	if fs.Exists("/opt/homebrew/bin/brew") {
		return true
	}

	if fs.Exists("/usr/local/bin/brew") {
		return true
	}

	return false
}

func (b *Brew) Install() error {

	err := command.Run(
		"bash",
		"-c",
		`NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
	)
	if err != nil {
		return err
	}

	return command.Run(
		"bash",
		"-c",
		`eval "$(/opt/homebrew/bin/brew shellenv)"`,
	)
}

func init() {
	Register(&Brew{})
}
