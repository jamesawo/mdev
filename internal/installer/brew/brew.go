package brew

import (
	"os/exec"

	"github.com/jamesawo/mdev/internal/runner"
)

func Install(pkg string) error {

	r := &runner.CommandRunner{}

	return r.Run("brew", "install", pkg)
}

func Uninstall(pkg string) error {

	r := &runner.CommandRunner{}

	return r.Run("brew", "uninstall", pkg)
}

func InstallCask(pkg string) error {

	r := &runner.CommandRunner{}

	return r.Run("brew", "install", "--cask", pkg)
}

func UninstallCask(pkg string) error {

	r := &runner.CommandRunner{}

	return r.Run("brew", "uninstall", "--cask", pkg)
}

func IsInstalled(pkg string) bool {

	cmd := exec.Command("brew", "list", pkg)

	err := cmd.Run()

	return err == nil
}
