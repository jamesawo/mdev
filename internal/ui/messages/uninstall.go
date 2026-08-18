package messages

import "fmt"

const (
	UninstallCmdShortDescription = "Uninstall a tool from your local environment"
	UninstallCmdLongDescription  = `Remove an installed tool from the local environment.

This command uninstalls a tool that was previously installed using
mdev. It removes the tool binaries and any managed directories that
belong to the mdev environment while keeping unrelated user files
untouched.

The command validates that the environment is configured and that
the specified tool is known by mdev before attempting removal.

Usage:
  mdev uninstall [tool]

Arguments:
  tool    Name of the tool to uninstall.

Behavior:
  • Verifies the mdev environment configuration.
  • Confirms the tool is supported by mdev.
  • Removes the installed tool from the managed environment.

Examples:
  mdev uninstall java
  mdev uninstall gradle
  mdev uninstall maven

Notes:
  Only tools managed by mdev can be removed using this command.
  If the tool is not installed, the command will exit without
  making changes.`
)

const (
	UninstallEnvironmentNotConfigured = "Environment not configured. Run `mdev doctor` first."
	UninstallDependencyWarning        = "Dependency warning"
	UninstallPlan                     = "Uninstall plan"
	UninstallDirectoriesToRemove      = "Directories to be removed"
	UninstallTools                    = "Uninstalling tools"
	UninstallRemoveDependentsQuestion = "Remove dependent tools first?"
	UninstallContinueQuestion         = "Continue uninstall?"
	UninstallCancelled                = "Cancelled."
)

func UninstallNotInstalled(name string) string {
	return fmt.Sprintf("%s not installed", name)
}

func UninstallRequiredBy(name string) string {
	return fmt.Sprintf("%s is required by:", name)
}

func UninstallRemoving(name string) string {
	return "Removing " + name
}

func UninstallCleaningStorage(path string) string {
	return "Cleaning storage: " + path
}

func UninstallRemoved(name string) string {
	return name + " removed"
}
