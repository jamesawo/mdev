package shell

import (
	"context"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
)

// RunWithSDKMAN executes a command after loading SDKMAN into the shell.
func RunWithSDKMAN(cmd string) error {
	return RunWithSDKMANContext(context.Background(), cmd)
}

// RunWithSDKMANContext executes an SDKMAN-backed command with cancellation.
func RunWithSDKMANContext(ctx context.Context, cmd string) error {
	bash, err := prerequisites.ModernBashPath(ctx)
	if err != nil {
		return err
	}
	return command.RunContext(
		ctx,
		bash,
		"-c",
		"source $HOME/.sdkman/bin/sdkman-init.sh && "+cmd,
	)
}
