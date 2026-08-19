package messages

import "fmt"

const (
	SetupCmdShortDescription = "configure where mdev stores managed data"
	SetupCmdLongDescription  = `Configure where mdev stores managed data.

Setup creates the selected mdev storage directory and saves its canonical path
in ~/.mdev/config.yaml. It does not install tools or move existing data.`
	SetupStoragePathFlag = "storage location or final mdev directory"
)

const (
	SetupStoragePathRequired = "configuration storage_path is required"
	SetupFailed              = "setup failed"
	SetupTitle               = "welcome to mdev."
	SetupCommand             = "mdev setup"
	SetupChooseDirectory     = "where should mdev store data?"
	SetupCustomStorageOption = "choose another location"
	SetupEnterStoragePath    = "storage path"
	SetupExistingTitle       = "an mdev folder already exists here."
	SetupUseExisting         = "use existing folder"
	SetupChooseAnother       = "choose another location"
	SetupCancel              = "cancel"
	SetupSymlinkTitle        = "this location is a symlink."
	SetupContinue            = "continue"
	SetupStorageUnavailable  = "storage is not available."
	SetupConnectRetry        = "connect drive and try again"
	SetupReady               = "mdev is ready."
	SetupNextStep            = "see what's available:"
	SetupCreateIt            = "create it"
	SetupTryAgain            = "try again"
)

func SetupDefaultStorageOption(location string) string {
	return fmt.Sprintf("%s (recommended)", displayHomePath(location))
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

func displayHomePath(path string) string { return path }
