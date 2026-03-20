package messages

import "fmt"

// --- Shared / general ---

const (
	System = "System"

	DoctorFailed = "doctor failed"
	DoctorReport = "Doctor Report"

	Mising     = "missing"
	Aborted    = "aborted"
	Installed  = "installed"
	Installing = "Installing"

	SetupCancelled  = "Setup cancelled"
	NoDriveDetected = "no drive was detected"
)

// --- Root command ---

const (
	RootAvailableCommands  = "Available commands"
	RootTypicalWorkflow    = "Typical workflow"
	RootCmdDoctor          = "mdev doctor   Initialize and validate your environment"
	RootCmdInstall         = "mdev install  Install development tools"
	RootCmdList            = "mdev list     Show supported tools and their status"
	RootCmdGraph           = "mdev graph    Show dependency graph between tools"
	RootCmdVersion         = "mdev version  Show version information"
	RootWorkflowDoctor     = "mdev doctor"
	RootWorkflowInstall    = "mdev install"
	RootWorkflowInstallAll = "mdev install --all"
)

// --- Version command ---

func VersionInfo(ver string) string    { return "mdev " + ver }
func VersionAuthor(author string) string { return "Created by " + author }

// --- Doctor command ---

const (
	DoctorTools                      = "Tools"
	DoctorNextSteps                  = "Next steps"
	DoctorInstallIndividual          = "Install individual tools:"
	DoctorInstallEverything          = "Install everything:"
	DoctorFixHint                    = "To fix system issues automatically:"
	DoctorFixCmd                     = "mdev doctor --fix"
	DoctorHintRunInstall             = "Run `mdev install`"
	DoctorNothingToFix               = "Nothing to fix"
	DoctorFixingPrerequisites        = "Fixing system prerequisites"
	DoctorInstallMissingPrerequisites = "Install missing prerequisites?"
	DoctorFixingSystem               = "Fixing system"
	DoctorFixSummary                 = "Summary"
)

func DoctorNotConfigured(name string) string {
	return fmt.Sprintf("%s not configured", name)
}

func DoctorInstallTool(name string) string {
	return fmt.Sprintf("mdev install %s", name)
}

func DoctorInstallationFailed(name string) string {
	return fmt.Sprintf("%s installation failed", name)
}

func DoctorTimeElapsed(d string) string {
	return "time: " + d
}

func DoctorTotalFixTime(d string) string {
	return "Total fix time: " + d
}

// --- Environment ---

const (
	Environment                        = "Environment"
	EnvironmentSetup                   = "Environment setup"
	EnvironmentChooseDirectory         = "Choose where to store development tool data."
	EnvironmentCreateDirectoryQuestion = "Create the directory now?"
	EnvironmentNotConfigured           = "Development data location not configured"
	EnvironmentNoDirectorySelected     = "No location selected, setup cancelled"
	EnvironmentSetupCompleted          = "Location setup done"
	EnvironmentSetupFailed             = "Location setup failed"
	EnvironmentLocation                = "Location"
)

// --- Install command ---

const (
	ErrInstallAllWithTool       = "cannot use --all with a specific tool"
	ErrEnvironmentNotConfigured = "Environment not configured. Run `mdev doctor` first."

	ToolsSelectToInstall         = "Select tools to install"
	ToolsInstallPlan             = "Install plan"
	ToolsContinueInstallQuestion = "Continue installation?"
	ToolsInstallCancelled        = "Installation cancelled."
	ToolsInstallingStart         = "Installing tools"
	ToolsAlreadyInstalled        = "already installed"

	ListInstalledSuffix = " (installed)"
)

func UnknownTool(name string) string {
	return fmt.Sprintf("unknown tool: %s", name)
}

// --- Uninstall command ---

const (
	UninstallDependencyWarning        = "Dependency warning"
	UninstallRemoveDependentsQuestion = "Remove dependent tools first?"
	UninstallCancelled                = "Cancelled."
	UninstallPlan                     = "Uninstall plan"
	UninstallDirectoriesToRemove      = "Directories to be removed"
	UninstallContinueQuestion         = "Continue uninstall?"
	UninstallingTools                 = "Uninstalling tools"
)

func UninstallRequiredBy(name string) string {
	return fmt.Sprintf("%s is required by:", name)
}

func UninstallNotInstalled(name string) string {
	return fmt.Sprintf("%s not installed", name)
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

// --- List command ---

const (
	ListAvailableTools = "Available tools"
)

// --- Graph command ---

const (
	GraphTitle = "Tool dependency graph"
)

// --- Interactive select ---

const (
	SelectEnterChoice  = "Enter the number of your choice, or 'q' to quit:"
	SelectPromptSymbol = ">"
	SelectInvalidInput = "Invalid input. Please enter a number."
	SelectOutOfRange   = "Selection out of range."
)

// --- Services ---

func ServiceFailedToStart(name string, err error) error {
	return fmt.Errorf("%s failed to start: %w", name, err)
}

func ServiceFailedToStop(name string, err error) error {
	return fmt.Errorf("%s failed to stop: %w", name, err)
}
