package sdkman

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

type SDKMAN struct{}

func (s *SDKMAN) StorageDir(env *environment.Environment) string {
	return storage.ToolDir(env, "sdkman")
}

func (s *SDKMAN) Name() string {
	return "sdkman"
}

func (s *SDKMAN) Description() string {
	return "Java version manager"
}

func (s *SDKMAN) IsInstalled(env *environment.Environment) bool {

	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	path := filepath.Join(home, ".sdkman")

	return fs.Exists(path)
}

func (s *SDKMAN) Install(env *environment.Environment) error {

	return command.Run("bash", "-c", "curl -s https://get.sdkman.io | bash")
}

func (s *SDKMAN) Configure(env *environment.Environment) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".sdkman")
	target := s.StorageDir(env)

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

func (s *SDKMAN) Verify(env *environment.Environment) error {

	cmd := exec.Command("bash", "-c", "source $HOME/.sdkman/bin/sdkman-init.sh && sdk version")

	return cmd.Run()
}

func (s *SDKMAN) Uninstall(env *environment.Environment) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".sdkman")

	if fs.IsSymlink(source) {
		return fs.Remove(source)
	}

	return nil
}

func (s *SDKMAN) Dependencies() []string {
	return []string{}
}

func init() {
	tools.Register(&SDKMAN{})
}
