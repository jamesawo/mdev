package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jamesawo/mdev/internal/ui/messages"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StoragePath string `yaml:"storage_path"`
}

var ErrStoragePathRequired = errors.New(messages.SetupStoragePathRequired)
var ErrAlreadyConfigured = errors.New("mdev is already configured")

func UserHomeDir() (string, error) {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" && sudoUser != "root" {
		account, err := user.Lookup(sudoUser)
		if err != nil {
			return "", fmt.Errorf("look up invoking user %q: %w", sudoUser, err)
		}
		if account.HomeDir == "" {
			return "", fmt.Errorf("invoking user %q has no home directory", sudoUser)
		}
		return account.HomeDir, nil
	}
	return os.UserHomeDir()
}

func configDir() string {
	home, err := UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".mdev")
}

func configFile() string {
	return filepath.Join(configDir(), "config.yaml")
}

func Exists() bool {
	_, err := os.Stat(configFile())
	return err == nil
}

func Save(cfg Config) error {
	if cfg.StoragePath == "" {
		return ErrStoragePathRequired
	}

	dir := configDir()
	if dir == "" {
		return fmt.Errorf("resolve configuration directory")
	}
	_, statErr := os.Stat(dir)
	dirCreated := errors.Is(statErr, os.ErrNotExist)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	if dirCreated {
		defer func() { _ = os.Remove(dir) }()
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	temporary, err := os.CreateTemp(dir, ".config-*.yaml")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := temporary.Chmod(0644); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Link(temporaryPath, configFile()); err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrAlreadyConfigured
		}
		return err
	}
	return nil
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configFile())
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.StoragePath == "" {
		return nil, ErrStoragePathRequired
	}

	return &cfg, nil
}
