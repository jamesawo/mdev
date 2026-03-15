package gradle

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/fs"
	brew "github.com/jamesawo/mdev/internal/infrastructure/packagemanager"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
	"github.com/jamesawo/mdev/internal/tools"
)

type Gradle struct{}

func (g *Gradle) StorageDir(env *environment.Environment) string {
	return storage.ToolDir(env, "gradle")
}

func (g *Gradle) Name() string {
	return "gradle"
}

func (g *Gradle) Description() string {
	return "Build automation tool"
}

func (g *Gradle) IsInstalled(env *environment.Environment) bool {
	return brew.IsInstalled("gradle")
}

func (g *Gradle) Install(env *environment.Environment) error {
	return brew.Install("gradle")
}

func (g *Gradle) Configure(env *environment.Environment) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".gradle")
	target := g.StorageDir(env)

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

func (g *Gradle) Verify(env *environment.Environment) error {

	cmd := exec.Command("gradle", "-version")

	return cmd.Run()
}

func (g *Gradle) Uninstall(env *environment.Environment) error {

	err := brew.Uninstall("gradle")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	source := filepath.Join(home, ".gradle")

	if fs.IsSymlink(source) {
		return fs.Remove(source)
	}

	return nil
}

func (g *Gradle) Dependencies() []string {
	return []string{"java"}
}

// register Gradle as a tool
func init() {
	tools.Register(&Gradle{})
}
