package prerequisites

import (
	"context"
	"errors"
	"os"
	"os/exec"
)

type Brew struct{}

func (Brew) Name() string {
	return "brew"
}

func (Brew) Check() bool {
	installed, _ := (Brew{}).InstallationStatus()
	return installed
}

func (Brew) InstallationStatus() (bool, error) {
	_, err := exec.LookPath("brew")
	if err == nil {
		return true, nil
	}
	if errors.Is(err, exec.ErrNotFound) {
		return false, nil
	}
	return false, err
}

func (Brew) Readiness(ctx context.Context) (State, string, error) {
	path, err := exec.LookPath("brew")
	if errors.Is(err, exec.ErrNotFound) {
		return StateMissing, "Homebrew is required", nil
	}
	if err != nil {
		return StateBroken, "", err
	}
	if output, err := exec.CommandContext(ctx, path, "--version").CombinedOutput(); err != nil {
		return StateBroken, string(output), nil
	}
	return StateReady, path, nil
}

func (Brew) Install() error {
	return (Brew{}).RemediateContext(context.Background())
}

func (Brew) PrerequisiteDependencies() []string { return []string{"curl", "xcode-cli"} }
func (Brew) RemediationDescription() string     { return "install Homebrew" }
func (Brew) RemediateContext(ctx context.Context) error {
	if err := runInstaller(ctx); err != nil {
		return err
	}

	refreshPath()

	return nil
}

func (Brew) VerifyContext(context.Context) error {
	_, err := exec.LookPath("brew")
	return err
}

// runInstaller executes the official Homebrew installation script
// and attaches the process to the user's terminal so interactive
// prompts (like sudo password) work correctly.
func runInstaller(ctx context.Context) error {

	cmd := exec.CommandContext(
		ctx,
		"/bin/bash",
		"-c",
		`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// refreshPath ensures the Homebrew binary directory is added
// to the current process PATH so brew can be discovered
// without restarting the shell.
func refreshPath() {

	brewPath := "/opt/homebrew/bin"

	current := os.Getenv("PATH")

	os.Setenv("PATH", brewPath+":"+current)
}

func init() {
	Register(Brew{})
}
