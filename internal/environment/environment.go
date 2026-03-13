package environment

import (
	"os"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/config"
)

type Environment struct {
	ExternalDrive string
	DataRoot      string
}

func New(externalDrive string) *Environment {
	return &Environment{
		ExternalDrive: externalDrive,
		DataRoot:      filepath.Join(externalDrive, "data"),
	}
}

func CreateDataRoot(env *Environment) error {
	return os.MkdirAll(env.DataRoot, 0755)
}

func FromConfig() (*Environment, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	env := New(cfg.ExternalDrive)

	return env, nil
}
