package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	StoragePath string `yaml:"storage_path"`
}

var ErrStoragePathRequired = errors.New("configuration storage_path is required")

func configDir() string {
	home, err := os.UserHomeDir()
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

	err := os.MkdirAll(configDir(), 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile(), data, 0644)
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
