package nvm

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/fs"
	"github.com/jamesawo/mdev/internal/storage"
	"github.com/jamesawo/mdev/internal/tools"
)

type NVM struct{}

func (n *NVM) Name() string {
	return "nvm"
}

func (n *NVM) Description() string {
	return "Node version manager"
}

func (n *NVM) IsInstalled(env *environment.Environment) bool {

	_, err := exec.LookPath("nvm")

	return err == nil
}

func (n *NVM) Install(env *environment.Environment) error {

	return command.Run(
		"bash",
		"-c",
		"curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash",
	)
}

func (n *NVM) Configure(env *environment.Environment) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".nvm")
	target := storage.ToolDir(env, "nvm")

	err = fs.EnsureDir(target)
	if err != nil {
		return err
	}

	if fs.IsSymlink(source) {
		return nil
	}

	if fs.Exists(source) {
		err = fs.Move(source, target)
		if err != nil {
			return err
		}
	}

	return fs.CreateSymlink(target, source)
}

func (n *NVM) Verify(env *environment.Environment) error {

	cmd := exec.Command(
		"bash",
		"-c",
		"source $HOME/.nvm/nvm.sh && nvm --version",
	)

	return cmd.Run()
}

func (n *NVM) Uninstall(env *environment.Environment) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".nvm")

	if fs.IsSymlink(source) {
		return fs.Remove(source)
	}

	return nil
}

func init() {
	tools.Register(&NVM{})
}
