package uninstall

import (
	"path/filepath"

	"github.com/jamesawo/mdev/internal/environment"
)

// StoragePath returns the directory used by a tool inside mdev storage.
func StoragePath(env *environment.Environment, name string) string {
	return filepath.Join(env.DataRoot, name)
}
