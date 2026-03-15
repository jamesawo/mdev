package shell

import "github.com/jamesawo/mdev/internal/command"

func RunWithSDKMAN(cmd string) error {

	return command.Run(
		"bash",
		"-c",
		"source $HOME/.sdkman/bin/sdkman-init.sh && "+cmd,
	)
}
