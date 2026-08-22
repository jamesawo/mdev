package nvm

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

type NVM struct{}

// Tool contract methods below describe NVM metadata, authoritative state,
// cancellable lifecycle operations, managed storage, and uninstall behavior.
func (*NVM) Name() string                                   { return "nvm" }
func (*NVM) Description() string                            { return messages.ToolsNVMDescription }
func (*NVM) Dependencies() []string                         { return nil }
func (*NVM) SystemPrerequisites() []string                  { return []string{"curl", "git"} }
func (*NVM) StorageDir(env *environment.Environment) string { return storage.ToolDir(env, "nvm") }
func (n *NVM) IsInstalled(env *environment.Environment) bool {
	installed, _ := n.InstallationStatus(env)
	return installed
}
func (n *NVM) InstallationStatus(env *environment.Environment) (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	managed, err := tools.ManagedSymlinkStatus(filepath.Join(home, ".nvm"), n.StorageDir(env))
	if err != nil || !managed {
		return managed, err
	}
	return tools.CommandInstallationStatus("bash", "-c", "source $HOME/.nvm/nvm.sh && nvm --version")
}
func (n *NVM) Install(env *environment.Environment) error {
	return n.InstallContext(context.Background(), env)
}
func (*NVM) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return command.RunContext(ctx, "bash", "-c", "curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash")
}
func (n *NVM) Configure(env *environment.Environment) error {
	return n.ConfigureContext(context.Background(), env)
}
func (n *NVM) ConfigureContext(ctx context.Context, env *environment.Environment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return storage.Relocate(filepath.Join(home, ".nvm"), n.StorageDir(env))
}
func (n *NVM) Verify(env *environment.Environment) error {
	return n.VerifyContext(context.Background(), env)
}
func (*NVM) VerifyContext(ctx context.Context, _ *environment.Environment) error {
	return shell.RunWithNVMContext(ctx, "nvm --version")
}
func (*NVM) Uninstall(_ *environment.Environment) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	source := filepath.Join(home, ".nvm")
	if info, err := os.Lstat(source); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return os.Remove(source)
	}
	return nil
}

func init() { tools.Register(&NVM{}) }
