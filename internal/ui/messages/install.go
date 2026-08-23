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
	InstallContinueQuestion    = "continue installation?"
	InstallCancelled           = "installation cancelled."
	InstallNoSelection         = "no tools selected."
	InstallPlan                = "install plan"
	InstallPlanItem            = "  %s\n"
	InstallProgress            = "%s %s...\n"
	InstallCompleted           = "✓ %s installed\n"
	InstallAlreadyInstalled    = "%s is already installed.\n"
	InstallUninstallHint       = "uninstall: mdev uninstall %s\n"
	InstallPhaseInstall        = "installing"
	InstallPhaseConfigure      = "configuring"
	InstallPhaseVerify         = "verifying"
	InstallActionInstall       = "install"
	InstallActionConfigure     = "configure"
	InstallActionVerify        = "verify"
	InstallNotConfigured       = "mdev is not configured; run mdev setup"
	InstallConfigurationError  = "mdev configuration cannot be read: %w"
	InstallStorageUnavailable  = "configured storage is unavailable at %s: %v"
	InstallStorageNotDirectory = "configured storage is unavailable at %s: expected a directory"
	InstallStorageNotWritable  = "configured storage is not writable at %s: %v"
	InstallUnknownTool         = "unknown tool %q. run `mdev list` to see available tools."
	InstallPlanError           = "resolve install plan: %w"
	InstallStatusError         = "determine %s installation status: %w"
	InstallPhaseError          = "%s %s: %w"
	InstallReadinessError      = "check system readiness: %w"
	InstallPrerequisiteMissing = "system prerequisite %s is %s; rerun mdev setup or use mdev doctor for diagnosis"
)
