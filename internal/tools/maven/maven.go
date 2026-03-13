package maven

import (
	"os/exec"

	"github.com/jamesawo/mdev/internal/environment"
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
	return nil
}

func (m *Maven) Configure(env *environment.Environment) error {
	return nil
}

func (m *Maven) Verify(env *environment.Environment) error {
	return nil
}

func (m *Maven) Uninstall(env *environment.Environment) error {
	return nil
}

func init() {
	tools.Register(&Maven{})
}
