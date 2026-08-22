package sdkman

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/fs"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type SDKMAN struct{}

func (s *SDKMAN) StorageDir(env *environment.Environment) string {
	return storage.ToolDir(env, "sdkman")
}

func (s *SDKMAN) Name() string {
	return "sdkman"
}

func (s *SDKMAN) Description() string {
	return messages.ToolsSDKMANDescription
}

func (s *SDKMAN) IsInstalled(env *environment.Environment) bool {
	installed, _ := s.InstallationStatus(env)
	return installed
}

func (s *SDKMAN) InstallationStatus(_ *environment.Environment) (bool, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	path := filepath.Join(home, ".sdkman")
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
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
