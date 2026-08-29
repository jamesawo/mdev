package messages

const (
	UninstallCommandUse          = "uninstall [tool]"
	UninstallCmdShortDescription = "uninstall a managed development tool"
	UninstallCmdLongDescription  = `remove a registered tool and its managed data.

mdev resolves installed dependants, presents the complete removal plan, asks
for confirmation, and removes approved tools in safe reverse dependency order.
provider output is captured while mdev reports truthful lifecycle progress.`

	UninstallEnvironmentError         = "mdev configuration cannot be read: %w"
	UninstallDependencyWarning        = "dependency warning"
	UninstallPlan                     = "uninstall plan"
	UninstallDirectoriesToRemove      = "directories to be removed"
	UninstallRemoveDependentsQuestion = "remove dependent tools first?"
	UninstallContinueQuestion         = "continue uninstall?"
	UninstallCancelled                = "uninstall cancelled."
	UninstallItem                     = "  %s\n"
	UninstallDependantItem            = "    %s\n"
	UninstallToolProgress             = "uninstalling %s... "
	UninstallStorageProgress          = "cleaning %s storage... "
	UninstallPhaseSucceeded           = "✓\n"
	UninstallPhaseFailed              = "✗\n"
	UninstallPhaseCancelled           = "cancelled\n"
	UninstallRemoved                  = "%s removed.\n"
	UninstallNotInstalled             = "%s is not installed.\n"
	UninstallRequiredBy               = "  %s is required by:\n"
	UninstallPhaseTool                = "uninstall"
	UninstallPhaseStorage             = "storage"
	UninstallPlanError                = "resolve uninstall plan: %w"
	UninstallToolError                = "uninstall %s: %w"
	UninstallStorageError             = "clean %s storage: %w"
)
