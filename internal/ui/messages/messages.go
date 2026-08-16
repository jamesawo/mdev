package messages

import "fmt"

// Cmd text
const (
	CmdDoctorShortDescription = "Inspect system, environment, and tools"
	CmdDoctorLongDescription  = `Analyze your system and development environment.

This command reports missing prerequisites, environment issues,
and tool installation status.

Use --fix to attempt automatic remediation.`
	CmdGraphShortDescription = "Show the dependency graph of supported tools"
	CmdGraphLongDescription  = `Graph displays the dependency relationships between development
tools managed by mdev.

Some tools depend on others to function correctly. For example,
Java requires SDKMAN to install and manage versions, while tools
like Maven and Gradle require Java to be present before they can
be installed.

This command prints the dependency graph so you can understand the
order in which tools will be installed when running:

  mdev install --all

The output shows each tool and the tools it depends on.

Example:

  mdev graph

Possible output:

  java -> sdkman
  maven -> java
  gradle -> java
  nvm
  podman
`
	CmdInstallShortDescription = "Install a tool in your local environment."
	CmdInstallLongDescription  = `
	Install a development tool into your local environment.

This command installs a supported tool and prepares it for use with the
current mdev environment configuration. The tool will be downloaded,
installed, and configured using the paths and settings defined in your
mdev environment.

Before running this command, your environment must be initialized and
validated using 'mdev doctor'. The install process depends on the
configured directories, tool paths, and system checks performed during
that step.

If the tool is already installed, the command will detect it and skip
the installation to avoid overwriting an existing setup.

Usage:
  mdev install [tool]

Arguments:
  tool    Name of the tool to install.

Behavior:
  • Validates that the environment is configured.
  • Checks whether the requested tool is supported.
  • Detects if the tool is already installed.
  • Runs the tool-specific installation process.

Examples:
  mdev install java
  mdev install gradle
  mdev install maven

Notes:
  Each tool provides its own installation logic. The command acts as a
  dispatcher that resolves the requested tool and executes its install
  routine using the current environment configuration.
	`
	CmdListShortDescription = "List all supported development tools"
	CmdListLongDescription  = `List all development tools supported by mdev.

This command displays the tools that mdev knows how to manage.
Each tool includes a short description and can be installed,
configured, and managed through the mdev lifecycle.

Typical usage:

  mdev list
`
	CmdRootShortDescription = "Automate development environment setup on macOS"
	CmdRootLongDescription  = `mdev is a command-line tool for setting up and managing
a development environment on macOS.

It installs development tools, configures them, and relocates large
tool caches to external storage to keep your system disk clean.`
	CmdSetupShortDescription = "Configure the development environment"
	CmdSetupLongDescription  = `Configure where mdev stores development tool data.

This command guides you through selecting an external drive, creates the
managed data directory, and saves the environment configuration for use by
other mdev commands.`
	CmdUninstallShortDescription = "Uninstall a tool from your local environment"
	CmdUninstallLongDescription  = `Remove an installed tool from the local environment.

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
	CmdUpShortDescription = "Start development services (like podman, ollama)"
	CmdUpLongDescription  = `Start the development environment.

mdev up starts the runtime services and tools required during development.

This typically includes:
- podman machine
- ollama service


Use this command at the beginning of your development session.`
	CmdVersionShortDescription = "Show mdev version information"
	CmdVersionLongDescription  = `Display the current version of mdev and basic
information about the project.`
)

const (
	FlagDoctorFix  = "Attempt to fix detected issues"
	FlagConfirmAll = "Accept all confirmation prompts"
)

// Section headers
const (
	System                       = "System"
	DoctorTools                  = "Tools"
	DoctorNextSteps              = "Next steps"
	DoctorFixingPrerequisites    = "Fixing system prerequisites"
	DoctorFixingSystem           = "Fixing system"
	DoctorFixSummary             = "Summary"
	Environment                  = "Environment"
	EnvironmentSetup             = "Environment setup"
	EnvironmentConfiguration     = "Environment configuration"
	ToolsInstallPlan             = "Install plan"
	ToolsInstallingStart         = "Installing tools"
	UninstallDependencyWarning   = "Dependency warning"
	UninstallPlan                = "Uninstall plan"
	UninstallDirectoriesToRemove = "Directories to be removed"
	UninstallingTools            = "Uninstalling tools"
	ListAvailableTools           = "Available tools"
	GraphTitle                   = "Tool dependency graph"
)

// Errors
const (
	DoctorFailed                   = "doctor failed"
	ErrInstallAllWithTool          = "cannot use --all with a specific tool"
	ErrEnvironmentNotConfigured    = "Environment not configured. Run `mdev doctor` first."
	EnvironmentNoDirectorySelected = "No location selected, setup cancelled"
	EnvironmentSetupFailed         = "Location setup failed"
	EnvironmentNotConfiguredShort  = "Not configured"
	EnvironmentNotSetUp            = "Development environment is not set up"
)

// Prompts
const (
	EnvironmentCreateDirectoryQuestion = "Create the directory now?"
	ToolsContinueInstallQuestion       = "Continue installation?"
	UninstallRemoveDependentsQuestion  = "Remove dependent tools first?"
	UninstallContinueQuestion          = "Continue uninstall?"
	DoctorInstallMissingPrerequisites  = "Install missing prerequisites?"
	SelectEnterChoice                  = "Enter the number of your choice, or 'q' to quit:"
	SelectPromptSymbol                 = ">"
	SelectInvalidInput                 = "Invalid input. Please enter a number."
	SelectOutOfRange                   = "Selection out of range."
)

func EnvironmentUseInternalStorage(location string) string {
	return fmt.Sprintf("No external drive detected. Use %s instead?", location)
}

// Info and status
const (
	Run                        = "run"
	Missing                    = "missing"
	Aborted                    = "aborted"
	Installed                  = "installed"
	Installing                 = "Installing"
	EnvironmentSetupCompleted  = "Location setup done"
	EnvironmentSetupCommand    = "mdev setup"
	EnvironmentLocation        = "Location"
	EnvironmentChooseDirectory = "Choose where to store development tool data."
	ToolsSelectToInstall       = "Select tools to install"
	ToolsInstallCancelled      = "Installation cancelled."
	ToolsAlreadyInstalled      = "already installed"
	ListInstalledSuffix        = " (installed)"
	UninstallCancelled         = "Cancelled."
	IndentMark                 = "└─ "
	StorageLocation            = "mdev directory"
	DataDirectory              = "data directory"
)

const (
	DoctorNothingToFix      = "Nothing to fix"
	DoctorInstallIndividual = "Install individual tools:"
	DoctorInstallEverything = "Install everything:"
	DoctorFixHint           = "To fix system issues automatically:"
)

// Tool descriptions

const (
	ToolJavaDescription   = "Java runtime (via SDKMAN)"
	ToolNVMDescription    = "Node version manager"
	ToolSDKMANDescription = "Java version manager"
	ToolMavenDescription  = "Java build automation tool"
	ToolGradleDescription = "Build automation tool"
	ToolPodmanDescription = "Container runtime with Podman Desktop"
)

func DoctorInstallationFailed(name string) string {
	return fmt.Sprintf("%s installation failed", name)
}

func UninstallNotInstalled(name string) string {
	return fmt.Sprintf("%s not installed", name)
}

func UnknownTool(name string) string {
	return fmt.Sprintf("unknown tool: %s", name)
}

func ServiceFailedToStart(name string, err error) error {
	return fmt.Errorf("%s failed to start: %w", name, err)
}

func ServiceFailedToStop(name string, err error) error {
	return fmt.Errorf("%s failed to stop: %w", name, err)
}

func DoctorInstallTool(name string) string {
	return fmt.Sprintf("mdev install %s", name)
}

func DoctorTimeElapsed(d string) string {
	return "time: " + d
}

func DoctorTotalFixTime(d string) string {
	return "Total fix time: " + d
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

func VersionInfo(ver string) string      { return "mdev " + ver }
func VersionAuthor(author string) string { return "Created by " + author }
