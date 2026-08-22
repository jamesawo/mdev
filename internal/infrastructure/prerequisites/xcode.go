package prerequisites

import (
	"errors"
	"os/exec"

	"github.com/jamesawo/mdev/internal/command"
)

type Xcode struct{}

func (x *Xcode) Name() string {
	return "xcode-cli"
}

func (x *Xcode) Check() bool {
	return true
}

func (x *Xcode) InstallationStatus() (bool, error) {
	err := exec.Command("xcode-select", "-p").Run()
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return false, nil
	}
	var commandErr *exec.Error
	if errors.As(err, &commandErr) && errors.Is(commandErr.Err, exec.ErrNotFound) {
		return false, nil
	}
	return false, err
}

func (x *Xcode) Install() error {
	return command.Run("xcode-select", "--install")
}

func init() {
	Register(&Xcode{})
}
