package podmandesktop

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	brew "github.com/jamesawo/mdev/internal/infrastructure/packagemanager"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

const caskName = "podman-desktop"

type PodmanDesktop struct{}

// Tool contract methods below keep the optional Desktop application separate
// from the Podman CLI and managed machine owned by the podman tool.
func (*PodmanDesktop) Name() string                               { return caskName }
func (*PodmanDesktop) Description() string                        { return messages.ToolsPodmanDesktopDescription }
func (*PodmanDesktop) Dependencies() []string                     { return []string{"podman"} }
func (*PodmanDesktop) SystemPrerequisites() []string              { return []string{"brew"} }
func (*PodmanDesktop) StorageDir(*environment.Environment) string { return "" }
func (p *PodmanDesktop) IsInstalled(env *environment.Environment) bool {
	return p.InstallationStatus(env)
}

// InstallationStatus recognizes both a Homebrew-owned cask and an existing
// user-managed application so install never overwrites an application it does
// not own.
func (*PodmanDesktop) InstallationStatus(*environment.Environment) bool {
	if brew.IsCaskInstalled(caskName) {
		return true
	}
	for _, path := range desktopApplicationPaths() {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (p *PodmanDesktop) Install(env *environment.Environment) error {
	return p.InstallContext(context.Background(), env)
}

func (*PodmanDesktop) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return brew.InstallCaskContext(ctx, caskName)
}

func (*PodmanDesktop) Configure(*environment.Environment) error { return nil }

func (*PodmanDesktop) ConfigureContext(ctx context.Context, _ *environment.Environment) error {
	return ctx.Err()
}

func (p *PodmanDesktop) Verify(env *environment.Environment) error {
	return p.VerifyContext(context.Background(), env)
}

func (p *PodmanDesktop) VerifyContext(ctx context.Context, env *environment.Environment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !p.InstallationStatus(env) {
		return errors.New(messages.ToolsPodmanDesktopNotInstalled)
	}
	return nil
}

// Uninstall removes only a cask Homebrew identifies as installed. An existing
// application outside that ownership boundary is deliberately preserved.
func (*PodmanDesktop) Uninstall(*environment.Environment) error {
	if !brew.IsCaskInstalled(caskName) {
		return nil
	}
	return brew.UninstallCask(caskName)
}

func desktopApplicationPaths() []string {
	paths := []string{filepath.Join(string(filepath.Separator), "Applications", "Podman Desktop.app")}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, "Applications", "Podman Desktop.app"))
	}
	return paths
}

func init() { tools.Register(&PodmanDesktop{}) }
