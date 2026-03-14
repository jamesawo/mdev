package podman

import (
	"os/exec"

	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/installer/brew"
	"github.com/jamesawo/mdev/internal/tools"
)

type Podman struct{}

func (p *Podman) Name() string {
	return "podman"
}

func (p *Podman) Description() string {
	return "Container runtime with Podman Desktop"
}

func (p *Podman) IsInstalled(env *environment.Environment) bool {
	return brew.IsInstalled("podman-desktop")
}

func (p *Podman) Install(env *environment.Environment) error {
	return brew.InstallCask("podman-desktop")
}

func (p *Podman) Configure(env *environment.Environment) error {
	return nil
}

func (p *Podman) Verify(env *environment.Environment) error {

	cmd := exec.Command("podman", "--version")

	return cmd.Run()
}

func (p *Podman) Uninstall(env *environment.Environment) error {
	return brew.UninstallCask("podman-desktop")
}

func init() {
	tools.Register(&Podman{})
}
