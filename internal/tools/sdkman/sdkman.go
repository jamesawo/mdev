package sdkman

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/shell"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type SDKMAN struct{}

// Tool contract methods below describe SDKMAN metadata, authoritative state,
// cancellable lifecycle operations, managed storage, and uninstall behavior.
func (*SDKMAN) Name() string                                   { return "sdkman" }
func (*SDKMAN) Description() string                            { return messages.ToolsSDKMANDescription }
func (*SDKMAN) Dependencies() []string                         { return nil }
func (*SDKMAN) StorageDir(env *environment.Environment) string { return storage.ToolDir(env, "sdkman") }
func (s *SDKMAN) IsInstalled(env *environment.Environment) bool {
	installed, _ := s.InstallationStatus(env)
	return installed
}
func (s *SDKMAN) InstallationStatus(env *environment.Environment) (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	managed, err := tools.ManagedSymlinkStatus(filepath.Join(home, ".sdkman"), s.StorageDir(env))
	if err != nil || !managed {
		return managed, err
	}
	return tools.CommandInstallationStatus("bash", "-c", "source $HOME/.sdkman/bin/sdkman-init.sh && sdk version")
}
func (s *SDKMAN) Install(env *environment.Environment) error {
	return s.InstallContext(context.Background(), env)
}
func (*SDKMAN) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return command.RunContext(ctx, "bash", "-c", "curl -s https://get.sdkman.io | bash")
}
func (s *SDKMAN) Configure(env *environment.Environment) error {
	return s.ConfigureContext(context.Background(), env)
}
func (s *SDKMAN) ConfigureContext(ctx context.Context, env *environment.Environment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return storage.Relocate(filepath.Join(home, ".sdkman"), s.StorageDir(env))
}
func (s *SDKMAN) Verify(env *environment.Environment) error {
	return s.VerifyContext(context.Background(), env)
}
func (*SDKMAN) VerifyContext(ctx context.Context, _ *environment.Environment) error {
	return shell.RunWithSDKMANContext(ctx, "sdk version")
}
func (*SDKMAN) Uninstall(_ *environment.Environment) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	source := filepath.Join(home, ".sdkman")
	if info, err := os.Lstat(source); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return os.Remove(source)
	}
	return nil
}

func init() { tools.Register(&SDKMAN{}) }
