package prerequisites

import (
	"os"
	"os/exec"
)

type Brew struct{}

func (Brew) Name() string {
	return "brew"
}

func (Brew) Check() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (Brew) Install() error {

	if err := runInstaller(); err != nil {
		return err
	}

	refreshPath()

	return nil
}

// runInstaller executes the official Homebrew installation script
// and attaches the process to the user's terminal so interactive
// prompts (like sudo password) work correctly.
func runInstaller() error {

	cmd := exec.Command(
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

	os.Setenv("PATH", current+":"+brewPath)
}

func init() {
	Register(Brew{})
}
