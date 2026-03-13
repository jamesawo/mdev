package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ExternalDrive string `yaml:"external_drive"`
}

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

	return &cfg, nil
}
