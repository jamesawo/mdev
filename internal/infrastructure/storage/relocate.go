package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Relocate conservatively moves a tool's conventional state directory into
// managed storage and replaces the source with a validated symlink.
func Relocate(source, target string) error {
	symlink, err := isSymlink(source)
	if err != nil {
		return err
	}
	if symlink {
		return reconcileSymlink(source, target)
	}

	sourceInfo, sourceErr := os.Stat(source)
	targetInfo, targetErr := os.Stat(target)
	if sourceErr != nil && !errors.Is(sourceErr, os.ErrNotExist) {
		return fmt.Errorf("inspect source %s: %w", source, sourceErr)
	}
	if targetErr != nil && !errors.Is(targetErr, os.ErrNotExist) {
		return fmt.Errorf("inspect target %s: %w", target, targetErr)
	}
	sourceExists := sourceErr == nil
	targetExists := targetErr == nil
	if sourceExists && !sourceInfo.IsDir() {
		return fmt.Errorf("tool state source is not a directory: %s", source)
	}
	if targetExists && !targetInfo.IsDir() {
		return fmt.Errorf("managed tool target is not a directory: %s", target)
	}

	if sourceExists && targetExists {
		empty, err := directoryEmpty(target)
		if err != nil {
			return err
		}
		if !empty {
			return fmt.Errorf("tool state exists at both %s and %s; move or remove one location and retry", source, target)
		}
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("remove empty managed target %s: %w", target, err)
		}
		targetExists = false
	}

	if sourceExists {
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		if err := os.Rename(source, target); err != nil {
			return fmt.Errorf("move tool state from %s to %s: %w", source, target, err)
		}
	} else if !targetExists {
		if err := os.MkdirAll(target, 0755); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(source), 0755); err != nil {
		return fmt.Errorf("prepare tool state parent for %s: %w", source, err)
	}
	if err := os.Symlink(target, source); err != nil {
		return fmt.Errorf("link tool state from %s to %s: %w", source, target, err)
	}
	return nil
}

// ManagedSymlinkStatus reports whether source points to an existing managed target.
// A correctly targeted dangling link is incomplete state, not an inspection error.
func ManagedSymlinkStatus(source, target string) (bool, error) {
	symlink, err := isSymlink(source)
	if err != nil {
		return false, err
	}
	if !symlink {
		return false, nil
	}
	matches, err := symlinkPointsTo(source, target)
	if err != nil {
		return false, err
	}
	if !matches {
		return false, nil
	}
	info, err := os.Stat(target)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect managed tool target %s: %w", target, err)
	}
	return info.IsDir(), nil
}

// reconcileSymlink accepts only the intended link and recreates its missing target.
func reconcileSymlink(source, target string) error {
	matches, err := symlinkPointsTo(source, target)
	if err != nil {
		return err
	}
	if !matches {
		destination, readErr := os.Readlink(source)
		if readErr != nil {
			return fmt.Errorf("inspect tool state symlink %s: %w", source, readErr)
		}
		return fmt.Errorf("tool state symlink %s points to %s, expected %s", source, destination, target)
	}
	info, err := os.Stat(target)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(target, 0755); err != nil {
			return fmt.Errorf("prepare managed tool target %s: %w", target, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect managed tool target %s: %w", target, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("managed tool target is not a directory: %s", target)
	}
	return nil
}

func symlinkPointsTo(source, target string) (bool, error) {
	destination, err := os.Readlink(source)
	if err != nil {
		return false, fmt.Errorf("inspect tool state symlink %s: %w", source, err)
	}
	if !filepath.IsAbs(destination) {
		destination = filepath.Join(filepath.Dir(source), destination)
	}
	destination, err = filepath.Abs(destination)
	if err != nil {
		return false, err
	}
	expected, err := filepath.Abs(target)
	if err != nil {
		return false, err
	}
	return filepath.Clean(destination) == filepath.Clean(expected), nil
}

// isSymlink distinguishes missing paths from existing symbolic links.
func isSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

// directoryEmpty determines whether a managed target can be safely replaced.
func directoryEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("inspect managed target %s: %w", path, err)
	}
	return len(entries) == 0, nil
}
