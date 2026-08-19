package environment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

type Volume struct {
	Name     string
	Path     string
	Writable bool
}

type locationChoice struct {
	path   string
	custom bool
}

var saveConfig = config.Save

func SetupInteractive() (*Environment, error) {
	if env, configured, err := Existing(); err != nil {
		return nil, err
	} else if configured {
		if info, statErr := os.Stat(env.StoragePath); statErr == nil && info.IsDir() {
			return env, ErrAlreadyConfigured
		} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return nil, statErr
		}
		return waitForConfiguredStorage(env)
	}

	defaultLocation, err := DefaultLocation()
	if err != nil {
		return nil, err
	}
	for {
		choice, err := chooseLocation(defaultLocation)
		if err != nil {
			return nil, err
		}
		if choice.custom {
			inputPath, pathErr := absoluteInputPath(choice.path)
			if pathErr != nil {
				return nil, pathErr
			}
			if _, statErr := os.Stat(inputPath); errors.Is(statErr, os.ErrNotExist) {
				printer.Section("that location doesn't exist.")
				printer.Info(inputPath)
				missingChoice, err := selectChoice([]string{messages.SetupCreateIt, messages.SetupTryAgain, messages.SetupCancel})
				if err != nil {
					return nil, err
				}
				if missingChoice == 1 {
					continue
				}
				if missingChoice != 0 {
					return nil, ErrSetupCancelled
				}
			}
		}
		resolved, symlinkSource, err := ResolveStoragePath(choice.path)
		if err != nil {
			printer.Fail(err.Error())
			if retry, retryErr := retryLocation(); retryErr != nil {
				return nil, retryErr
			} else if retry {
				continue
			}
			return nil, ErrSetupCancelled
		}

		if symlinkSource != "" {
			printer.Section(messages.SetupSymlinkTitle)
			printer.Info(symlinkSource + "  → " + resolved)
			printer.Info(messages.SetupWillUse(resolved))
			choice, err := selectChoice([]string{messages.SetupContinue, messages.SetupChooseAnother, messages.SetupCancel})
			if err != nil {
				return nil, err
			}
			if choice == 1 {
				continue
			}
			if choice != 0 {
				return nil, ErrSetupCancelled
			}
		}

		if info, statErr := os.Stat(resolved); statErr == nil && info.IsDir() {
			printer.Section(messages.SetupExistingTitle)
			choice, err := selectChoice([]string{messages.SetupUseExisting, messages.SetupChooseAnother, messages.SetupCancel})
			if err != nil {
				return nil, err
			}
			if choice == 1 {
				continue
			}
			if choice != 0 {
				return nil, ErrSetupCancelled
			}
		}

		env, err := SetupResolved(resolved)
		if err == nil {
			return env, nil
		}
		printer.Fail(err.Error())
		if retry, retryErr := retryLocation(); retryErr != nil {
			return nil, retryErr
		} else if retry {
			continue
		}
		return nil, ErrSetupCancelled
	}
}

func chooseLocation(defaultLocation string) (locationChoice, error) {
	volumes, err := ListExternalVolumes()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return locationChoice{}, err
	}
	options := []string{messages.SetupDefaultStorageOption(DisplayPath(defaultLocation))}
	for _, volume := range volumes {
		options = append(options, messages.SetupVolumeOption(volume.Name, volume.Path, volume.Writable))
	}
	options = append(options, messages.SetupCustomStorageOption)

	for {
		index, err := interactive.RadioSelect(messages.SetupChooseDirectory, options)
		if err != nil {
			return locationChoice{}, interruptionError(err)
		}
		switch {
		case index < 0:
			return locationChoice{}, ErrSetupCancelled
		case index == 0:
			return locationChoice{path: defaultLocation}, nil
		case index == len(options)-1:
			location, err := interactive.Input(messages.SetupEnterStoragePath, defaultLocation)
			if err != nil {
				return locationChoice{}, interruptionError(err)
			}
			return locationChoice{path: location, custom: true}, nil
		default:
			volume := volumes[index-1]
			if !volume.Writable {
				printer.Fail(messages.SetupReadOnly(volume.Path))
				continue
			}
			return locationChoice{path: volume.Path}, nil
		}
	}
}

func retryLocation() (bool, error) {
	choice, err := selectChoice([]string{messages.SetupChooseAnother, messages.SetupCancel})
	if err != nil {
		return false, err
	}
	return choice == 0, nil
}

func waitForConfiguredStorage(env *Environment) (*Environment, error) {
	for {
		printer.Section(messages.SetupStorageUnavailable)
		printer.Info(messages.SetupExpected(env.StoragePath))
		choice, err := selectChoice([]string{messages.SetupConnectRetry, messages.SetupCancel})
		if err != nil {
			return nil, err
		}
		if choice != 0 {
			return nil, ErrSetupCancelled
		}
		if info, statErr := os.Stat(env.StoragePath); statErr == nil && info.IsDir() {
			return env, ErrAlreadyConfigured
		}
	}
}

func selectChoice(options []string) (int, error) {
	choice, err := interactive.RadioSelect("", options)
	if err != nil {
		return -1, interruptionError(err)
	}
	return choice, nil
}

func interruptionError(err error) error {
	if errors.Is(err, terminal.InterruptErr) {
		return ErrSetupCancelled
	}
	return err
}

func DefaultLocation() (string, error) {
	home, err := config.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "mdev"), nil
}

// DisplayPath shortens paths below the invoking user's home for calm,
// user-facing output while persisted paths remain canonical and absolute.
func DisplayPath(path string) string {
	home, err := config.UserHomeDir()
	if err != nil {
		return path
	}
	canonicalHome, _, err := resolveExistingPrefix(home)
	if err != nil {
		canonicalHome = filepath.Clean(home)
	}
	canonicalPath := filepath.Clean(path)
	relative, err := filepath.Rel(canonicalHome, canonicalPath)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return path
	}
	if relative == "." {
		return "~"
	}
	return filepath.Join("~", relative)
}

// ResolveStoragePath turns a selected parent or mdev directory into the
// canonical absolute mdev-owned storage path. Environment variables are not
// expanded, and symlinks in the longest existing prefix are resolved.
func ResolveStoragePath(path string) (resolved string, symlinkSource string, err error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", "", fmt.Errorf("storage path must not be empty")
	}
	home, err := config.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	absPath, err := absoluteInputPath(path)
	if err != nil {
		return "", "", err
	}
	if absPath == string(filepath.Separator) || absPath == filepath.Clean(home) {
		return "", "", fmt.Errorf("storage location must be a dedicated directory, not %s", absPath)
	}
	if filepath.Base(absPath) != "mdev" {
		absPath = filepath.Join(absPath, "mdev")
	}

	physical, usedSymlink, err := resolveExistingPrefix(absPath)
	if err != nil {
		return "", "", err
	}
	if usedSymlink {
		symlinkSource = filepath.Clean(path)
	}
	info, statErr := os.Stat(physical)
	if statErr == nil && !info.IsDir() {
		return "", "", fmt.Errorf("location is a file, not a directory: %s", physical)
	}
	if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
		return "", "", fmt.Errorf("inspect storage location: %w", statErr)
	}
	return physical, symlinkSource, nil
}

func absoluteInputPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	home, err := config.UserHomeDir()
	if err != nil {
		return "", err
	}
	if path == "~" {
		path = home
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(absPath), nil
}

func ValidateStoragePath(path string) (string, error) {
	resolved, _, err := ResolveStoragePath(path)
	return resolved, err
}

func resolveExistingPrefix(path string) (string, bool, error) {
	prefix := path
	var suffix []string
	for {
		_, err := os.Lstat(prefix)
		if err == nil {
			break
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", false, err
		}
		parent := filepath.Dir(prefix)
		if parent == prefix {
			return "", false, err
		}
		suffix = append(suffix, filepath.Base(prefix))
		prefix = parent
	}
	physicalPrefix, err := filepath.EvalSymlinks(prefix)
	if err != nil {
		return "", false, err
	}
	resolved := physicalPrefix
	for i := len(suffix) - 1; i >= 0; i-- {
		resolved = filepath.Join(resolved, suffix[i])
	}
	return filepath.Clean(resolved), filepath.Clean(resolved) != filepath.Clean(path), nil
}

func Existing() (*Environment, bool, error) {
	cfg, err := config.Load()
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("configuration cannot be read; fix or remove ~/.mdev/config.yaml and try again: %w", err)
	}
	return New(cfg.StoragePath), true, nil
}

func Setup(location string) (*Environment, error) {
	resolved, _, err := ResolveStoragePath(location)
	if err != nil {
		return nil, err
	}
	return SetupResolved(resolved)
}

func SetupResolved(location string) (*Environment, error) {
	canonicalLocation, _, err := resolveExistingPrefix(location)
	if err != nil {
		return nil, err
	}
	location = canonicalLocation
	if err := validateResolvedStoragePath(location); err != nil {
		return nil, err
	}
	if _, configured, err := Existing(); err != nil {
		return nil, err
	} else if configured {
		return nil, ErrAlreadyConfigured
	}
	created, err := createStorageRoot(location)
	if err != nil {
		return nil, actionablePathError(location, err)
	}
	committed := false
	defer func() {
		if !committed {
			rollbackCreatedDirectories(created)
		}
	}()
	if err := saveConfig(config.Config{StoragePath: location}); err != nil {
		if errors.Is(err, config.ErrAlreadyConfigured) {
			return nil, ErrAlreadyConfigured
		}
		return nil, fmt.Errorf("save configuration: %w", err)
	}
	committed = true
	return New(location), nil
}

func validateResolvedStoragePath(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("storage path must be a clean absolute path")
	}
	if filepath.Base(path) != "mdev" {
		return fmt.Errorf("storage path must end in mdev")
	}
	home, err := config.UserHomeDir()
	if err != nil {
		return err
	}
	if path == string(filepath.Separator) || path == filepath.Clean(home) {
		return fmt.Errorf("storage path is too broad: %s", path)
	}
	return nil
}

func createStorageRoot(path string) ([]string, error) {
	var missing []string
	current := path
	for {
		_, err := os.Stat(current)
		if err == nil {
			break
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		missing = append(missing, current)
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	probe, err := os.CreateTemp(path, ".mdev-write-check-*")
	if err != nil {
		rollbackCreatedDirectories(missing)
		return nil, err
	}
	probePath := probe.Name()
	if err := probe.Close(); err != nil {
		_ = os.Remove(probePath)
		rollbackCreatedDirectories(missing)
		return nil, err
	}
	if err := os.Remove(probePath); err != nil {
		rollbackCreatedDirectories(missing)
		return nil, err
	}
	return missing, nil
}

func rollbackCreatedDirectories(paths []string) {
	for _, path := range paths {
		_ = os.Remove(path)
	}
}

func actionablePathError(path string, err error) error {
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("%s is not writable; choose another location or rerun setup with sudo: %w", path, err)
	}
	return fmt.Errorf("prepare storage at %s: %w", path, err)
}
