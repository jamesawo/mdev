package java

import (
	"context"
	"os/exec"
	"strings"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
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
	installed, err := tools.CommandInstallationStatus("java", "-version")
	if err != nil && strings.Contains(err.Error(), "Unable to locate a Java Runtime") {
		return false, nil
	}
	if err != nil || !installed {
		return installed, err
	}
	bash, err := prerequisites.ModernBashPath(context.Background())
	if err != nil {
		return false, nil
	}
	return tools.CommandInstallationStatus(bash, "-c", "source $HOME/.sdkman/bin/sdkman-init.sh && test -n \"$JAVA_HOME\"")
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
	bash, err := prerequisites.ModernBashPath(ctx)
	if err != nil {
		return err
	}
	return exec.CommandContext(ctx, bash, "-c", "source $HOME/.sdkman/bin/sdkman-init.sh && test -n \"$JAVA_HOME\"").Run()
}
func (j *Java) Verify(env *environment.Environment) error {
	return j.VerifyContext(context.Background(), env)
}
func (*Java) VerifyContext(ctx context.Context, _ *environment.Environment) error {
	return shell.RunWithSDKMANContext(ctx, "java -version")
}
func (*Java) Uninstall(_ *environment.Environment) error {
	return shell.RunWithSDKMAN("sdk uninstall java 21.0.8-tem")
}

func init() { tools.Register(&Java{}) }
