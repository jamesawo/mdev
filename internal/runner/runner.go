package runner

import (
	"os"
	"os/exec"
)

type Runner interface {
	Run(name string, args ...string) error
}

type CommandRunner struct{}

func (r *CommandRunner) Run(name string, args ...string) error {

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
