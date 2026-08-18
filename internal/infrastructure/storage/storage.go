package storage

import (
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

func ToolDir(env *environment.Environment, tool string) string {
	return filepath.Join(env.StoragePath, tool)
}
