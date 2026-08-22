package messages

import "fmt"

const (
	SetupCommandName         = "setup"
	SetupStoragePathFlagName = "storage-path"
	SetupYesFlagName         = "yes"
	SetupCmdShortDescription = "configure where mdev stores managed data"
	SetupCmdLongDescription  = `Prepare mdev for normal use.

Setup validates system prerequisites, asks before making system changes, then
creates the selected storage directory and saves its canonical path in
~/.mdev/config.yaml. It does not install registered development tools.`
	SetupStoragePathFlag = "storage location or final mdev directory"
	SetupYesUnsupported  = "setup does not support --yes"
)

const (
	SetupStoragePathRequired     = "configuration storage_path is required"
	SetupFailed                  = "setup failed"
	SetupTitle                   = "welcome to mdev."
	SetupCommand                 = "mdev setup"
	SetupChooseDirectory         = "where should mdev store data?"
	SetupCustomStorageOption     = "choose another location"
	SetupEnterStoragePath        = "storage path"
	SetupMissingLocation         = "that location doesn't exist."
	SetupExistingTitle           = "an mdev folder already exists here."
	SetupUseExisting             = "use existing folder"
	SetupChooseAnother           = "choose another location"
	SetupCancel                  = "cancel"
	SetupSymlinkTitle            = "this location is a symlink."
	SetupContinue                = "continue"
	SetupStorageUnavailable      = "storage is not available."
	SetupConnectRetry            = "connect drive and try again"
	SetupReady                   = "mdev is ready."
	SetupNextStep                = "see what's available:"
	SetupCreateIt                = "create it"
	SetupTryAgain                = "try again"
	SetupCancelled               = "setup cancelled."
	SetupComplete                = "setup is complete."
	SetupListCommand             = "mdev list"
	SetupAlreadyConfigured       = "mdev is already configured at %s; setup will not replace it"
	SetupError                   = "%s: %w"
	SetupStorageFormat           = "storage: %s"
	SetupCancelledError          = "setup cancelled"
	SetupConfiguredError         = "mdev is already configured"
	SetupSymlinkFormat           = "%s  → %s"
	SetupStoragePathEmpty        = "storage path must not be empty"
	SetupDedicatedStorage        = "storage location must be a dedicated directory, not %s"
	SetupLocationIsFile          = "location is a file, not a directory: %s"
	SetupInspectStorage          = "inspect storage location: %w"
	SetupConfigurationUnreadable = "configuration cannot be read; fix or remove ~/.mdev/config.yaml and try again: %w"
	SetupSaveConfiguration       = "save configuration: %w"
	SetupStorageCleanAbsolute    = "storage path must be a clean absolute path"
	SetupStorageEndsInMdev       = "storage path must end in mdev"
	SetupStorageTooBroad         = "storage path is too broad: %s"
	SetupStorageNotWritable      = "%s is not writable; choose another location or rerun setup with sudo: %w"
	SetupPrepareStorage          = "prepare storage at %s: %w"
	SetupSetOwnership            = "set invoking-user ownership on %s: %w"
	SetupLookupInvokingUser      = "look up invoking user %q: %w"
	SetupInvokingUserNoHome      = "invoking user %q has no home directory"
	SetupParseInvokingUID        = "parse invoking user UID %q: %w"
	SetupParseInvokingGID        = "parse invoking user GID %q: %w"
	SetupResolveConfigDir        = "resolve configuration directory"
	SetupReadinessApply          = "apply these changes?"
	SetupReadinessDeclined       = "setup is incomplete; required system changes were declined"
	SetupReadinessNonInteractive = "setup requires system changes; rerun interactively to review and approve them"
	SetupReadinessChecking       = "checking %s...\n"
	SetupReadinessReady          = "✓ %s ready\n"
	SetupReadinessNeeds          = "  %s  %s\n"
	SetupReadinessRemediating    = "%s %s...\n"
	SetupReadinessVerified       = "✓ %s verified\n"
)

func SetupStorage(path string) string { return fmt.Sprintf(SetupStorageFormat, path) }
func SetupSymlinkResolution(source, resolved string) string {
	return fmt.Sprintf(SetupSymlinkFormat, source, resolved)
}

func SetupDefaultStorageOption(location string) string {
	return fmt.Sprintf("%s (recommended)", location)
}

func SetupVolumeOption(name, path string, writable bool) string {
	if !writable {
		return fmt.Sprintf("%s  %s (read-only)", name, path)
	}
	return fmt.Sprintf("%s  %s", name, path)
}

func SetupReadOnly(path string) string {
	return fmt.Sprintf("%s is read-only; choose another location", path)
}

func SetupWillUse(path string) string  { return "mdev will use:  " + path }
func SetupExpected(path string) string { return "expected: " + path }
