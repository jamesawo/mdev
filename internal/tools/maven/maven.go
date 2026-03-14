package maven

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/fs"
	"github.com/jamesawo/mdev/internal/installer/brew"
	"github.com/jamesawo/mdev/internal/storage"
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
	return brew.IsInstalled("maven")
}

func (m *Maven) Install(env *environment.Environment) error {
	return brew.Install("maven")
}

func (m *Maven) Configure(env *environment.Environment) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".m2")
	target := storage.ToolDir(env, "maven")

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

func (m *Maven) ConfigureOld(env *environment.Environment) error {

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

	err := brew.Uninstall("maven")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".m2")

	if fs.IsSymlink(source) {
		return fs.Remove(source)
	}

	return nil
}

// register Maven as a tool
func init() {
	tools.Register(&Maven{})
}
