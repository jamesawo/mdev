package shell

import (
	"context"

	"github.com/jamesawo/mdev/internal/command"
)

// RunWithNVM executes a command after loading NVM into the shell.
func RunWithNVM(cmd string) error {
	return RunWithNVMContext(context.Background(), cmd)
}

// RunWithNVMContext executes an NVM-backed command with caller cancellation.
func RunWithNVMContext(ctx context.Context, cmd string) error {
	return command.RunContext(
		ctx,
		"bash",
		"-c",
		"source $HOME/.nvm/nvm.sh && "+cmd,
	)
}
