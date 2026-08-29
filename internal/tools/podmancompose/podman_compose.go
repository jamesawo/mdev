package podmancompose

import (
	"context"
	"errors"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	brew "github.com/jamesawo/mdev/internal/infrastructure/packagemanager"
	"github.com/jamesawo/mdev/internal/infrastructure/subprocess"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

const formulaName = "podman-compose"

type PodmanCompose struct{}

// Tool contract methods below keep the optional Compose provider separate
// from the Podman CLI and managed machine owned by the podman tool.
func (*PodmanCompose) Name() string                               { return formulaName }
func (*PodmanCompose) Description() string                        { return messages.ToolsPodmanComposeDescription }
func (*PodmanCompose) Dependencies() []string                     { return []string{"podman"} }
func (*PodmanCompose) SystemPrerequisites() []string              { return []string{"brew"} }
func (*PodmanCompose) StorageDir(*environment.Environment) string { return "" }
func (p *PodmanCompose) IsInstalled(env *environment.Environment) bool {
	installed, _ := p.InstallationStatus(env)
	return installed
}

// InstallationStatus requires both Homebrew ownership and a working provider,
// so an unrelated executable does not count as mdev-managed completion.
func (*PodmanCompose) InstallationStatus(*environment.Environment) (bool, error) {
	if !brew.IsInstalled(formulaName) {
		return false, nil
	}
	return tools.CommandInstallationStatus(formulaName, "--version")
}

func (p *PodmanCompose) Install(env *environment.Environment) error {
	return p.InstallContext(context.Background(), env)
}

func (*PodmanCompose) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return brew.InstallContext(ctx, formulaName)
}

func (*PodmanCompose) Configure(*environment.Environment) error { return nil }

func (*PodmanCompose) ConfigureContext(ctx context.Context, _ *environment.Environment) error {
	return ctx.Err()
}

func (p *PodmanCompose) Verify(env *environment.Environment) error {
	return p.VerifyContext(context.Background(), env)
}

func (p *PodmanCompose) VerifyContext(ctx context.Context, env *environment.Environment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !brew.IsInstalled(formulaName) {
		return errors.New(messages.ToolsPodmanComposeNotManaged)
	}
	return subprocess.Check(ctx, formulaName, "--version")
}

// Uninstall removes only the Homebrew formula owned by this tool and leaves
// the Podman machine and other Compose providers untouched.
func (*PodmanCompose) Uninstall(*environment.Environment) error {
	if !brew.IsInstalled(formulaName) {
		return nil
	}
	return brew.Uninstall(formulaName)
}

func init() { tools.Register(&PodmanCompose{}) }
