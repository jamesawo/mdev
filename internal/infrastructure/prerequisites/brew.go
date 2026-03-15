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

	cmd := exec.Command(
		"/bin/bash",
		"-c",
		"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)",
	)

	// Attach to terminal so installer can ask for password
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func init() {
	Register(Brew{})
}
