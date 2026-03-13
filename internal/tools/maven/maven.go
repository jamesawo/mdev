package maven

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/runner"
	"github.com/jamesawo/mdev/internal/tools"
)

type Maven struct{}

func (m *Maven) Name() string {
	return "maven"
}

func (m *Maven) Description() string {
	return "Java build automation tool"
}

func (m *Maven) IsInstalled(env *environment.Environment) bool {

	_, err := exec.LookPath("mvn")

	return err == nil
}

func (m *Maven) Install(env *environment.Environment) error {

	r := &runner.CommandRunner{}

	return r.Run("brew", "install", "maven")
}

func (m *Maven) Configure(env *environment.Environment) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".m2")
	target := filepath.Join(env.DataRoot, "maven")

	err = os.MkdirAll(target, 0755)
	if err != nil {
		return err
	}

	info, err := os.Lstat(source)
	if err == nil {

		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		err = os.Rename(source, target)
		if err != nil {
			return err
		}
	}

	return os.Symlink(target, source)
}

func (m *Maven) Verify(env *environment.Environment) error {

	cmd := exec.Command("mvn", "-version")

	return cmd.Run()
}

func (m *Maven) Uninstall(env *environment.Environment) error {
	return nil
}

func init() {
	tools.Register(&Maven{})
}
