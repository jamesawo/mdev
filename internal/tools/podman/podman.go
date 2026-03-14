package podman

import (
	"os/exec"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/fs"
	"github.com/jamesawo/mdev/internal/installer/brew"
	"github.com/jamesawo/mdev/internal/storage"
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

	target := storage.ToolDir(env, "podman")

	err := fs.EnsureDir(target)
	if err != nil {
		return err
	}

	err = command.Run("podman", "machine", "init", "--image-path", target)
	if err != nil {
		return err
	}

	return command.Run("podman", "machine", "start")
}

func (p *Podman) Verify(env *environment.Environment) error {

	cmd := exec.Command("podman", "--version")

	return cmd.Run()
}

func (p *Podman) Uninstall(env *environment.Environment) error {
	return brew.UninstallCask("podman-desktop")
}

func (p *Podman) Dependencies() []string {
	return []string{}
}

func init() {
	tools.Register(&Podman{})
}
