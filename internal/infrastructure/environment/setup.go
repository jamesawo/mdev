package environment

import (
	"fmt"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/ui/interactive"
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
		return nil, fmt.Errorf("no drive was detected")
	}

	//printer.Section("Select a preferred drive for development data")

	index, err := interactive.RadioSelect(
		"Choose where to store development data",
		drives,
	)
	if err != nil {
		return nil, err
	}

	if index == -1 {
		return nil, fmt.Errorf("no location selected, Setup cancelled")
	}

	selected := drives[index]

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

	printer.Success("Location setup initialized")
	printer.Info("Location: " + externalDrive)

	return env, nil
}
