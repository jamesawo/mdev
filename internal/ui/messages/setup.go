package messages

import "fmt"

const (
	SetupCmdShortDescription = "Configure the development environment"
	SetupCmdLongDescription  = `Configure where mdev stores development tool data.

This command guides you through selecting a storage path, creates the managed
storage directory, and saves the environment configuration for use by other
mdev commands.`
)

const (
	SetupStoragePathRequired     = "configuration storage_path is required"
	SetupNoDirectorySelected     = "No location selected, setup cancelled"
	SetupFailed                  = "Location setup failed"
	SetupTitle                   = "Environment setup"
	SetupCreateDirectoryQuestion = "Create the directory now?"
	SetupCompleted               = "Location setup done"
	SetupCommand                 = "mdev setup"
	SetupLocation                = "Location"
	SetupChooseDirectory         = "Choose where to store development tool data."
	SetupCustomStorageOption     = "Choose another filesystem path"
	SetupEnterStoragePath        = "Storage path"
)

func SetupDefaultStorageOption(location string) string {
	return fmt.Sprintf("Default: %s", location)
}

func SetupUseStoragePath(location string) string {
	return fmt.Sprintf("Use %s as the mdev storage path?", location)
}
