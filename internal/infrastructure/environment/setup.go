package environment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type Volume struct {
	Name     string
	Path     string
	Writable bool
}

var (
	saveConfig = config.Save
	ownPath    = config.OwnPathForInvokingUser
)

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
		return "", "", errors.New(messages.SetupStoragePathEmpty)
	}
	home, err := config.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	absPath, err := AbsoluteInputPath(path)
	if err != nil {
		return "", "", err
	}
	if absPath == string(filepath.Separator) || absPath == filepath.Clean(home) {
		return "", "", fmt.Errorf(messages.SetupDedicatedStorage, absPath)
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
		return "", "", fmt.Errorf(messages.SetupLocationIsFile, physical)
	}
	if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
		return "", "", fmt.Errorf(messages.SetupInspectStorage, statErr)
	}
	return physical, symlinkSource, nil
}

// AbsoluteInputPath expands a leading home shorthand and resolves a clean
// absolute path without expanding environment variables.
func AbsoluteInputPath(path string) (string, error) {
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
		return nil, false, fmt.Errorf(messages.SetupConfigurationUnreadable, err)
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
		return nil, fmt.Errorf(messages.SetupSaveConfiguration, err)
	}
	committed = true
	return New(location), nil
}

func validateResolvedStoragePath(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return errors.New(messages.SetupStorageCleanAbsolute)
	}
	if filepath.Base(path) != "mdev" {
		return errors.New(messages.SetupStorageEndsInMdev)
	}
	home, err := config.UserHomeDir()
	if err != nil {
		return err
	}
	if path == string(filepath.Separator) || path == filepath.Clean(home) {
		return fmt.Errorf(messages.SetupStorageTooBroad, path)
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
	for _, createdPath := range missing {
		if err := ownPath(createdPath); err != nil {
			rollbackCreatedDirectories(missing)
			return nil, err
		}
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
		return fmt.Errorf(messages.SetupStorageNotWritable, path, err)
	}
	return fmt.Errorf(messages.SetupPrepareStorage, path, err)
}
