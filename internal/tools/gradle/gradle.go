package gradle

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	brew "github.com/jamesawo/mdev/internal/infrastructure/packagemanager"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type Gradle struct{}

// Tool contract methods below describe Gradle metadata, authoritative state,
// cancellable lifecycle operations, managed storage, and uninstall behavior.
func (*Gradle) Name() string                                   { return "gradle" }
func (*Gradle) Description() string                            { return messages.ToolsGradleDescription }
func (*Gradle) Dependencies() []string                         { return []string{"java"} }
func (*Gradle) SystemPrerequisites() []string                  { return []string{"brew"} }
func (*Gradle) StorageDir(env *environment.Environment) string { return storage.ToolDir(env, "gradle") }
func (g *Gradle) IsInstalled(env *environment.Environment) bool {
	installed, _ := g.InstallationStatus(env)
	return installed
}
func (g *Gradle) InstallationStatus(env *environment.Environment) (bool, error) {
	installed, err := tools.CommandInstallationStatus("gradle", "-version")
	if err != nil || !installed {
		return installed, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	return tools.ManagedSymlinkStatus(filepath.Join(home, ".gradle"), g.StorageDir(env))
}
func (g *Gradle) Install(env *environment.Environment) error {
	return g.InstallContext(context.Background(), env)
}
func (*Gradle) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return brew.InstallContext(ctx, "gradle")
}
func (g *Gradle) Configure(env *environment.Environment) error {
	return g.ConfigureContext(context.Background(), env)
}
func (g *Gradle) ConfigureContext(ctx context.Context, env *environment.Environment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return storage.Relocate(filepath.Join(home, ".gradle"), g.StorageDir(env))
}
func (g *Gradle) Verify(env *environment.Environment) error {
	return g.VerifyContext(context.Background(), env)
}
func (*Gradle) VerifyContext(ctx context.Context, _ *environment.Environment) error {
	return exec.CommandContext(ctx, "gradle", "-version").Run()
}
func (g *Gradle) Uninstall(env *environment.Environment) error {
	if err := brew.Uninstall("gradle"); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	source := filepath.Join(home, ".gradle")
	if info, err := os.Lstat(source); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return os.Remove(source)
	}
	return nil
}

func init() { tools.Register(&Gradle{}) }
