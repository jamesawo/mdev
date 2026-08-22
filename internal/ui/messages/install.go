package messages

const (
	InstallCommandUse          = "install [tool]"
	InstallAllFlagName         = "all"
	InstallAllFlag             = "install all registered tools"
	InstallCmdShortDescription = "install a development tool"
	InstallCmdLongDescription  = `make one or more registered development tools available.

each tool owns how it is installed, configured, verified, and stored. mdev
resolves tool dependencies and runs their lifecycle in dependency order.

usage:
  mdev install <tool>
  mdev install
  mdev install --all

install is progressive and retry-safe. it does not update valid installations,
start normal runtime state, manage services, or provide JSON output.`

	InstallAllWithToolError    = "cannot use --all with a specific tool"
	InstallSelectTools         = "select tools to install"
	InstallContinueQuestion    = "Continue installation?"
	InstallCancelled           = "Installation cancelled."
	InstallNoSelection         = "No tools selected."
	InstallPlan                = "Install plan"
	InstallPlanItem            = "  %s\n"
	InstallProgress            = "%s %s...\n"
	InstallCompleted           = "✓ %s installed\n"
	InstallAlreadyInstalled    = "%s is already installed.\n"
	InstallUninstallHint       = "Uninstall: mdev uninstall %s\n"
	InstallPhaseInstall        = "Installing"
	InstallPhaseConfigure      = "Configuring"
	InstallPhaseVerify         = "Verifying"
	InstallActionInstall       = "install"
	InstallActionConfigure     = "configure"
	InstallActionVerify        = "verify"
	InstallNotConfigured       = "mdev is not configured; run mdev setup"
	InstallConfigurationError  = "mdev configuration cannot be read: %w"
	InstallStorageUnavailable  = "configured storage is unavailable at %s: %v"
	InstallStorageNotDirectory = "configured storage is unavailable at %s: expected a directory"
	InstallStorageNotWritable  = "configured storage is not writable at %s: %v"
	InstallUnknownTool         = "Unknown tool %q. Run `mdev list` to see available tools."
	InstallPlanError           = "resolve install plan: %w"
	InstallStatusError         = "determine %s installation status: %w"
	InstallPhaseError          = "%s %s: %w"
)
