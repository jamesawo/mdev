package prerequisites

import (
	"context"
	"errors"
	"os/exec"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type Xcode struct{}

func (x *Xcode) Name() string {
	return "xcode-cli"
}

func (x *Xcode) Check() bool {
	installed, _ := x.InstallationStatus()
	return installed
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
	return x.RemediateContext(context.Background())
}

func (x *Xcode) RemediationDescription() string { return messages.ReadinessInstallXcodeCLI }
func (x *Xcode) RemediateContext(ctx context.Context) error {
	return command.RunContext(ctx, "xcode-select", "--install")
}
func (x *Xcode) VerifyContext(context.Context) error {
	installed, err := x.InstallationStatus()
	if err != nil {
		return err
	}
	if !installed {
		return errors.New(messages.ReadinessXcodeVerificationFailed)
	}
	return nil
}

func init() {
	Register(&Xcode{})
}
