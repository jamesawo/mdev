package packagemanager

import (
	"context"
	"os/exec"
)

// Install installs a Homebrew formula.
func Install(pkg string) error {
	return InstallContext(context.Background(), pkg)
}

// InstallContext installs a Homebrew formula with caller cancellation.
func InstallContext(ctx context.Context, pkg string) error {

	r := &CommandRunner{}

	return r.RunContext(ctx, "brew", "install", pkg)
}

func Uninstall(pkg string) error {

	r := &CommandRunner{}

	return r.Run("brew", "uninstall", pkg)
}

// InstallCask installs a Homebrew cask.
func InstallCask(pkg string) error {
	return InstallCaskContext(context.Background(), pkg)
}

// InstallCaskContext installs a Homebrew cask with caller cancellation.
func InstallCaskContext(ctx context.Context, pkg string) error {

	r := &CommandRunner{}

	return r.RunContext(ctx, "brew", "install", "--cask", pkg)
}

func UninstallCask(pkg string) error {

	r := &CommandRunner{}

	return r.Run("brew", "uninstall", "--cask", pkg)
}

func IsInstalled(pkg string) bool {

	cmd := exec.Command("brew", "list", pkg)

	err := cmd.Run()

	return err == nil
}
