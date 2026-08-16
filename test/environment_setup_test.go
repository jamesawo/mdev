package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

func TestDefaultLocationUsesMdevDirectoryInUserHome(t *testing.T) {
	home := givenUserHome(t)

	location, err := environment.DefaultLocation()

	assertDefaultLocationResolved(t, err)
	assertDefaultLocationIsInsideUserHome(t, location, home)
}

func givenUserHome(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)

	return home
}

func assertDefaultLocationResolved(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected default location to resolve: %v", err)
	}
}

func assertDefaultLocationIsInsideUserHome(t *testing.T, location string, home string) {
	t.Helper()

	expected := filepath.Join(home, "mdev")
	if location != expected {
		t.Fatalf("expected default location %q, got %q", expected, location)
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
