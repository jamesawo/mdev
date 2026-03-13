package environment

import (
	"os"
	"path/filepath"
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

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mdev")
}

func ConfigFile() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}
