package maven

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/shell"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type Maven struct{}

// Tool contract methods below describe Maven metadata, authoritative state,
// cancellable lifecycle operations, managed storage, and uninstall behavior.
func (*Maven) Name() string                                   { return "maven" }
func (*Maven) Description() string                            { return messages.ToolsMavenDescription }
func (*Maven) Dependencies() []string                         { return []string{"java"} }
func (*Maven) SystemPrerequisites() []string                  { return nil }
func (*Maven) StorageDir(env *environment.Environment) string { return storage.ToolDir(env, "maven") }
func (m *Maven) IsInstalled(env *environment.Environment) bool {
	installed, _ := m.InstallationStatus(env)
	return installed
}
func (m *Maven) InstallationStatus(env *environment.Environment) (bool, error) {
	installed, err := shell.SDKMANCandidateInstallationStatus(context.Background(), "maven", "mvn", "--version")
	if err != nil || !installed {
		return installed, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	return tools.ManagedSymlinkStatus(filepath.Join(home, ".m2"), m.StorageDir(env))
}
func (m *Maven) Install(env *environment.Environment) error {
	return m.InstallContext(context.Background(), env)
}
func (*Maven) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return shell.InstallSDKMANCandidateContext(ctx, "maven")
}
func (m *Maven) Configure(env *environment.Environment) error {
	return m.ConfigureContext(context.Background(), env)
}
func (m *Maven) ConfigureContext(ctx context.Context, env *environment.Environment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return storage.Relocate(filepath.Join(home, ".m2"), m.StorageDir(env))
}
func (m *Maven) Verify(env *environment.Environment) error {
	return m.VerifyContext(context.Background(), env)
}
func (*Maven) VerifyContext(ctx context.Context, _ *environment.Environment) error {
	return shell.RunSDKMANCandidateContext(ctx, "maven", "mvn", "--version")
}
func (m *Maven) Uninstall(env *environment.Environment) error {
	if err := shell.UninstallSDKMANCandidate("maven"); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	source := filepath.Join(home, ".m2")
	if info, err := os.Lstat(source); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return os.Remove(source)
	}
	return nil
}

func init() { tools.Register(&Maven{}) }
