package environment

import (
	"os"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
)

type Environment struct {
	StoragePath string
}

func New(storagePath string) *Environment {
	return &Environment{
		StoragePath: filepath.Clean(storagePath),
	}
}

func CreateStorageRoot(env *Environment) error {
	if err := os.MkdirAll(env.StoragePath, 0755); err != nil {
		return err
	}

	probe, err := os.CreateTemp(env.StoragePath, ".mdev-write-check-*")
	if err != nil {
		return err
	}
	probePath := probe.Name()

	if err := probe.Close(); err != nil {
		_ = os.Remove(probePath)
		return err
	}

	return os.Remove(probePath)
}

func FromConfig() (*Environment, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	env := New(cfg.StoragePath)

	return env, nil
}
