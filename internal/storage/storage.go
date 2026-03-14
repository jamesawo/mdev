package storage

import (
	"path/filepath"

	"github.com/jamesawo/mdev/internal/environment"
)

// todo: could we get the source and target from stroage?
func ToolDir(env *environment.Environment, tool string) string {
	return filepath.Join(env.DataRoot, tool)
}
