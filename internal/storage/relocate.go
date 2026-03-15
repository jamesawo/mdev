package storage

import "github.com/jamesawo/mdev/internal/fs"

// TODO(mdev-refactor):
// Tools currently duplicate the pattern of relocating a cache/config directory
// from the user's home (e.g. ~/.gradle, ~/.m2, ~/.nvm) into the mdev external
// storage and creating a symlink back.
//
// Future improvement:
//   1. Add SourceDir() to the Tool interface.
//   2. Use a shared helper like storage.Relocate(source, target).
//   3. Tools would only declare:
//        - SourceDir()
//        - StorageDir()
//      and the relocation logic would be centralized.
//
// This will reduce duplication across tool Configure() methods and make it
// easier to add new tools (Android SDK, IntelliJ, VSCode, etc.).

func Relocate(source string, target string) error {
	if err := fs.EnsureDir(target); err != nil {
		return err
	}

	if fs.IsSymlink(source) {
		return nil
	}

	if fs.Exists(source) {
		if err := fs.Move(source, target); err != nil {
			return err
		}
	}

	return fs.CreateSymlink(target, source)
}
