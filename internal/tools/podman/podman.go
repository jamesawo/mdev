package podman

import (
	"context"
	"os"
	"os/exec"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	brew "github.com/jamesawo/mdev/internal/infrastructure/packagemanager"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type Podman struct{}

// Tool contract methods below describe Podman metadata, authoritative state,
// cancellable provisioning, managed storage, and uninstall behavior.
func (*Podman) Name() string                                   { return "podman" }
func (*Podman) Description() string                            { return messages.ToolsPodmanDescription }
func (*Podman) Dependencies() []string                         { return nil }
func (*Podman) SystemPrerequisites() []string                  { return []string{"brew"} }
func (*Podman) StorageDir(env *environment.Environment) string { return storage.ToolDir(env, "podman") }
func (p *Podman) IsInstalled(env *environment.Environment) bool {
	installed, _ := p.InstallationStatus(env)
	return installed
}
func (p *Podman) InstallationStatus(env *environment.Environment) (bool, error) {
	installed, err := tools.CommandInstallationStatus("podman", "--version")
	if err != nil || !installed {
		return installed, err
	}
	if info, err := os.Stat(p.StorageDir(env)); err != nil || !info.IsDir() {
		return false, nil
	}
	if err := exec.Command("podman", "machine", "inspect").Run(); err != nil {
		return false, nil
	}
	return true, nil
}
func (p *Podman) Install(env *environment.Environment) error {
	return p.InstallContext(context.Background(), env)
}
func (*Podman) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return brew.InstallCaskContext(ctx, "podman-desktop")
}
func (p *Podman) Configure(env *environment.Environment) error {
	return p.ConfigureContext(context.Background(), env)
}
func (p *Podman) ConfigureContext(ctx context.Context, env *environment.Environment) error {
	target := p.StorageDir(env)
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	return command.RunContext(ctx, "podman", "machine", "init", "--image-path", target)
}
func (p *Podman) Verify(env *environment.Environment) error {
	return p.VerifyContext(context.Background(), env)
}
func (*Podman) VerifyContext(ctx context.Context, _ *environment.Environment) error {
	if err := exec.CommandContext(ctx, "podman", "--version").Run(); err != nil {
		return err
	}
	return exec.CommandContext(ctx, "podman", "machine", "inspect").Run()
}
func (*Podman) Uninstall(_ *environment.Environment) error {
	return brew.UninstallCask("podman-desktop")
}

func init() { tools.Register(&Podman{}) }
