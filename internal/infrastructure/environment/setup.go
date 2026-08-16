package environment

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// SetupInteractive initializes the environment by asking the user where
// development data should be stored.
func SetupInteractive() (*Environment, error) {

	drives, err := listExternalDrives()
	if err != nil {
		return nil, err
	}

	if len(drives) == 0 {
		location, err := DefaultLocation()
		if err != nil {
			return nil, err
		}

		if !confirmation.Ask(messages.EnvironmentUseInternalStorage(location)) {
			return nil, ErrSetupCancelled
		}

		return Setup(location)
	}

	index, err := interactive.RadioSelect(
		messages.EnvironmentChooseDirectory,
		drives,
	)
	if err != nil {
		return nil, err
	}

	if index == -1 {
		return nil, fmt.Errorf(messages.EnvironmentNoDirectorySelected)
	}

	return Setup(filepath.Join("/Volumes", drives[index]))
}

// DefaultLocation returns the standard internal storage location.
func DefaultLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "mdev"), nil
}

// Setup creates the managed data directory and saves its location.
// Existing directories and their contents are preserved.
func Setup(location string) (*Environment, error) {
	env := New(location)

	if err := CreateDataRoot(env); err != nil {
		return nil, err
	}

	if err := config.Save(config.Config{ExternalDrive: location}); err != nil {
		return nil, err
	}

	printer.Success(messages.EnvironmentSetupCompleted)
	printer.Info(messages.EnvironmentLocation + ": " + location)

	return env, nil
}
