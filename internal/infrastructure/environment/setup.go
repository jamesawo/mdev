package environment

import (
	"fmt"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// SetupInteractive initializes the environment by asking the user
// to choose which external drive should store development data.
func SetupInteractive() (*Environment, error) {

	drives, err := listExternalDrives()
	if err != nil {
		return nil, err
	}

	if len(drives) == 0 {
		return nil, fmt.Errorf(messages.NoDriveDetected)
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

	selected := drives[index]

	// todo: extract /volumns to a shared const file
	externalDrive := filepath.Join("/Volumes", selected)

	env := New(externalDrive)

	err = CreateDataRoot(env)
	if err != nil {
		return nil, err
	}

	cfg := config.Config{
		ExternalDrive: externalDrive,
	}

	err = config.Save(cfg)
	if err != nil {
		return nil, err
	}

	printer.Success(messages.EnvironmentSetupCompleted)
	printer.Info(messages.EnvironmentLocation + ": " + externalDrive)

	return env, nil
}
