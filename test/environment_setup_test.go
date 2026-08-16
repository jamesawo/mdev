package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

func TestDefaultLocationUsesMdevDirectoryInUserHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	location, err := environment.DefaultLocation()
	if err != nil {
		t.Fatalf("DefaultLocation() error = %v", err)
	}

	want := filepath.Join(home, "mdev")
	if location != want {
		t.Fatalf("DefaultLocation() = %q, want %q", location, want)
	}
}

func TestSetupPreservesExistingContents(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	location := filepath.Join(home, "mdev")
	existingFile := filepath.Join(location, "existing.txt")

	if err := os.MkdirAll(location, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(existingFile, []byte("keep me"), 0644); err != nil {
		t.Fatal(err)
	}

	env, err := environment.Setup(location)
	if err != nil {
		t.Fatalf("Setup() error = %v", err)
	}

	contents, err := os.ReadFile(existingFile)
	if err != nil {
		t.Fatalf("existing file was not preserved: %v", err)
	}
	if string(contents) != "keep me" {
		t.Fatalf("existing file contents = %q", contents)
	}
	if _, err := os.Stat(env.DataRoot); err != nil {
		t.Fatalf("data directory was not created: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}
	if cfg.ExternalDrive != location {
		t.Fatalf("configured location = %q, want %q", cfg.ExternalDrive, location)
	}
}
