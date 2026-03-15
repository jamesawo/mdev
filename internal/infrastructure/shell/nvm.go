package shell

import "github.com/jamesawo/mdev/internal/command"

func RunWithNVM(cmd string) error {

	return command.Run(
		"bash",
		"-c",
		"source $HOME/.nvm/nvm.sh && "+cmd,
	)
}
