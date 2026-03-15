package packagemanager

import (
	"os"
	"os/exec"
)

type CommandRunner struct{}

func (r *CommandRunner) Run(name string, args ...string) error {

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
