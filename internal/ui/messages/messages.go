package messages

import "fmt"

// Section headers

const (
	System                    = "System"
	DoctorReport              = "Doctor Report"
	DoctorTools               = "Tools"
	DoctorNextSteps           = "Next steps"
	DoctorFixingPrerequisites = "Fixing system prerequisites"
	DoctorFixingSystem        = "Fixing system"
	DoctorFixSummary          = "Summary"
	Environment               = "Environment"
	EnvironmentSetup          = "Environment setup"
	ToolsInstallPlan          = "Install plan"
	ToolsInstallingStart      = "Installing tools"
	UninstallDependencyWarning   = "Dependency warning"
	UninstallPlan                = "Uninstall plan"
	UninstallDirectoriesToRemove = "Directories to be removed"
	UninstallingTools            = "Uninstalling tools"
	ListAvailableTools           = "Available tools"
	GraphTitle                   = "Tool dependency graph"
	RootAvailableCommands        = "Available commands"
	RootTypicalWorkflow          = "Typical workflow"
)

// Errors

const (
	DoctorFailed                = "doctor failed"
	ErrInstallAllWithTool       = "cannot use --all with a specific tool"
	ErrEnvironmentNotConfigured = "Environment not configured. Run `mdev doctor` first."
	NoDriveDetected             = "no drive was detected"
	EnvironmentNotConfigured    = "Development data location not configured"
	EnvironmentNoDirectorySelected = "No location selected, setup cancelled"
	EnvironmentSetupFailed      = "Location setup failed"
)

func DoctorNotConfigured(name string) string {
	return fmt.Sprintf("%s not configured", name)
}

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

// Info and status

const (
	Missing    = "missing"
	Aborted                 = "aborted"
	Installed               = "installed"
	Installing              = "Installing"
	SetupCancelled          = "Setup cancelled"
	EnvironmentSetupCompleted = "Location setup done"
	EnvironmentLocation     = "Location"
	EnvironmentChooseDirectory = "Choose where to store development tool data."
	ToolsSelectToInstall    = "Select tools to install"
	ToolsInstallCancelled   = "Installation cancelled."
	ToolsAlreadyInstalled   = "already installed"
	ListInstalledSuffix     = " (installed)"
	UninstallCancelled      = "Cancelled."
	DoctorNothingToFix      = "Nothing to fix"
	DoctorHintRunInstall    = "Run `mdev install`"
	DoctorInstallIndividual = "Install individual tools:"
	DoctorInstallEverything = "Install everything:"
	DoctorFixHint           = "To fix system issues automatically:"
	DoctorFixCmd            = "mdev doctor --fix"
)

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

// Commands shown to users

const (
	RootCmdDoctor          = "mdev doctor   Initialize and validate your environment"
	RootCmdInstall         = "mdev install  Install development tools"
	RootCmdList            = "mdev list     Show supported tools and their status"
	RootCmdGraph           = "mdev graph    Show dependency graph between tools"
	RootCmdVersion         = "mdev version  Show version information"
	RootWorkflowDoctor     = "mdev doctor"
	RootWorkflowInstall    = "mdev install"
	RootWorkflowInstallAll = "mdev install --all"
)

func VersionInfo(ver string) string     { return "mdev " + ver }
func VersionAuthor(author string) string { return "Created by " + author }

// Tool descriptions

const (
	ToolJavaDescription   = "Java runtime (via SDKMAN)"
	ToolNVMDescription    = "Node version manager"
	ToolSDKMANDescription = "Java version manager"
	ToolMavenDescription  = "Java build automation tool"
	ToolGradleDescription = "Build automation tool"
	ToolPodmanDescription = "Container runtime with Podman Desktop"
)
