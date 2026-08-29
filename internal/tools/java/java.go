package java

import (
	"context"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/shell"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type Java struct{}

// Tool contract methods below describe Java metadata, authoritative state,
// cancellable SDKMAN lifecycle operations, and uninstall behavior.
func (*Java) Name() string                               { return "java" }
func (*Java) Description() string                        { return messages.ToolsJavaDescription }
func (*Java) Dependencies() []string                     { return []string{"sdkman"} }
func (*Java) StorageDir(*environment.Environment) string { return "" }
func (j *Java) IsInstalled(env *environment.Environment) bool {
	installed, _ := j.InstallationStatus(env)
	return installed
}
func (*Java) InstallationStatus(_ *environment.Environment) (bool, error) {
	return shell.SDKMANCandidateInstallationStatus(context.Background(), "java", "java", "-version")
}
func (j *Java) Install(env *environment.Environment) error {
	return j.InstallContext(context.Background(), env)
}
func (*Java) InstallContext(ctx context.Context, _ *environment.Environment) error {
	return shell.RunWithSDKMANContext(ctx, "sdk install java 21.0.8-tem")
}
func (j *Java) Configure(env *environment.Environment) error {
	return j.ConfigureContext(context.Background(), env)
}
func (*Java) ConfigureContext(ctx context.Context, _ *environment.Environment) error {
	return shell.RunSDKMANCandidateContext(ctx, "java", "java", "-version")
}
func (j *Java) Verify(env *environment.Environment) error {
	return j.VerifyContext(context.Background(), env)
}
func (*Java) VerifyContext(ctx context.Context, _ *environment.Environment) error {
	return shell.RunSDKMANCandidateContext(ctx, "java", "java", "-version")
}
func (*Java) Uninstall(_ *environment.Environment) error {
	return shell.UninstallSDKMANCandidateContext(context.Background(), "java")
}
func (*Java) UninstallContext(ctx context.Context, _ *environment.Environment) error {
	return shell.UninstallSDKMANCandidateContext(ctx, "java")
}

func init() { tools.Register(&Java{}) }
