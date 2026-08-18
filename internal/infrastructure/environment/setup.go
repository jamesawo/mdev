package environment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// SetupInteractive initializes the environment by asking the user where
// development data should be stored.
func SetupInteractive() (*Environment, error) {
	defaultLocation, err := DefaultLocation()
	if err != nil {
		return nil, err
	}

	volumes, err := listExternalVolumes()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	options := []string{messages.EnvironmentDefaultStorageOption(defaultLocation)}
	options = append(options, volumes...)
	options = append(options, messages.EnvironmentCustomStorageOption)

	index, err := interactive.RadioSelect(
		messages.EnvironmentChooseDirectory,
		options,
	)
	if err != nil {
		return nil, err
	}

	if index == -1 {
		return nil, fmt.Errorf(messages.EnvironmentNoDirectorySelected)
	}

	location := defaultLocation
	if index > 0 && index < len(options)-1 {
		location = volumes[index-1]
	}
	if index == len(options)-1 {
		location, err = interactive.Input(messages.EnvironmentEnterStoragePath, defaultLocation)
		if err != nil {
			return nil, err
		}
	}

	location, err = ValidateStoragePath(location)
	if err != nil {
		return nil, err
	}

	if !confirmation.Ask(messages.EnvironmentUseStoragePath(location)) {
		return nil, ErrSetupCancelled
	}

	return Setup(location)
}

// DefaultLocation returns the standard internal storage location.
func DefaultLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "mdev"), nil
}

// ValidateStoragePath normalizes a user-selected storage root and rejects
// locations that are too broad or cannot safely serve as directories.
func ValidateStoragePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("storage path must not be empty")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		path = home
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)

	if path == string(filepath.Separator) || path == filepath.Clean(home) {
		return "", fmt.Errorf("storage path must be a dedicated directory, not %s", path)
	}

	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return "", fmt.Errorf("storage path is not a directory: %s", path)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return path, nil
}

// Setup creates the managed storage directory and saves its location.
// Existing directories and their contents are preserved.
func Setup(location string) (*Environment, error) {
	location, err := ValidateStoragePath(location)
	if err != nil {
		return nil, err
	}

	env := New(location)

	if err := CreateStorageRoot(env); err != nil {
		return nil, err
	}

	if err := config.Save(config.Config{StoragePath: location}); err != nil {
		return nil, err
	}

	printer.Success(messages.EnvironmentSetupCompleted)
	printer.Info(messages.EnvironmentLocation + ": " + location)

	return env, nil
}
