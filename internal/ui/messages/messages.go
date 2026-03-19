package messages

// todo: put all the message across the ui into this constant for now
const (
	System = "System"

	DoctorFailed = "doctor failed"
	DoctorReport = "Doctor Report"

	Mising     = "missing"
	Aborted    = "aborted"
	Installed  = "installed"
	Installing = "Installing"

	Environment                        = "Environment"
	EnvironmentSetup                   = "Environment setup"
	EnvironmentChooseDirectory         = "Choose where to store development tool data."
	EnvironmentCreateDirectoryQuestion = "Create the directory now?"
	EnvironmentNotConfigured           = "Development data location not configured"
	EnvironmentNoDirectorySelected     = "No location selected, setup cancelled"
	EnvironmentSetupCompleted          = "Location setup done"
	EnvironmentSetupFailed             = "Location setup failed"
	EnvironmentLocation                = "Location"

	ToolsSelectToInstall         = "Select tools to install"
	ToolsInstallPlan             = "Install plan"
	ToolsContinueInstallQuestion = "Continue installation?"
	ToolsInstallCancelled        = "Installation cancelled."
	ToolsInstallingStart         = "Installing tools"
	ToolsAlreadyInstalled        = "already installed"

	SetupCancelled  = "Setup cancelled"
	NoDriveDetected = "no drive was detected"
)
