package environment

import (
	"fmt"
	"os"

	"github.com/jamesawo/mdev/internal/ui/messages"
)

// ValidateInstallStorage requires the configured storage root to already exist
// as a writable directory without leaving its write probe behind.
func ValidateInstallStorage(env *Environment) error {
	info, err := os.Stat(env.StoragePath)
	if err != nil {
		return fmt.Errorf(messages.InstallStorageUnavailable, env.StoragePath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf(messages.InstallStorageNotDirectory, env.StoragePath)
	}
	probe, err := os.CreateTemp(env.StoragePath, ".mdev-install-check-*")
	if err != nil {
		return fmt.Errorf(messages.InstallStorageNotWritable, env.StoragePath, err)
	}
	name := probe.Name()
	if err := probe.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf(messages.InstallStorageNotWritable, env.StoragePath, err)
	}
	if err := os.Remove(name); err != nil {
		return fmt.Errorf(messages.InstallStorageNotWritable, env.StoragePath, err)
	}
	return nil
}
