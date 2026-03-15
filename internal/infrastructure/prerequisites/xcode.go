package prerequisites

import "github.com/jamesawo/mdev/internal/command"

type Xcode struct{}

func (x *Xcode) Name() string {
	return "xcode-cli"
}

func (x *Xcode) Check() bool {
	return CommandExists("xcode-select")
}

func (x *Xcode) Install() error {
	return command.Run("xcode-select", "--install")
}

func init() {
	Register(&Xcode{})
}
