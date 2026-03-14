package java

import (
	"os/exec"

	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/tools"
)

type Java struct{}

func (j *Java) Name() string {
	return "java"
}

func (j *Java) Description() string {
	return "Java runtime (via SDKMAN)"
}

func (j *Java) IsInstalled(env *environment.Environment) bool {

	cmd := exec.Command("java", "-version")

	err := cmd.Run()

	return err == nil
}

func (j *Java) Install(env *environment.Environment) error {

	cmd := exec.Command(
		"bash",
		"-c",
		"source $HOME/.sdkman/bin/sdkman-init.sh && sdk install java 21.0.8-tem",
	)

	return cmd.Run()
}

func (j *Java) Configure(env *environment.Environment) error {

	cmd := exec.Command(
		"bash",
		"-c",
		"source $HOME/.sdkman/bin/sdkman-init.sh && echo $JAVA_HOME",
	)

	return cmd.Run()
}

func (j *Java) Verify(env *environment.Environment) error {

	cmd := exec.Command("java", "-version")

	return cmd.Run()
}

func (j *Java) Uninstall(env *environment.Environment) error {

	// todo: version 21 is the default for now, will be adjusted so user can change that later on
	cmd := exec.Command(
		"bash",
		"-c",
		"source $HOME/.sdkman/bin/sdkman-init.sh && sdk uninstall java 21.0.8-tem",
	)

	return cmd.Run()
}

func (j *Java) Dependencies() []string {
	return []string{"sdkman"}
}

func init() {
	tools.Register(&Java{})
}
