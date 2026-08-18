package test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/storage"
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
	if _, err := os.Stat(env.StoragePath); err != nil {
		t.Fatalf("storage directory was not created: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}
	if cfg.StoragePath != location {
		t.Fatalf("configured location = %q, want %q", cfg.StoragePath, location)
	}
}

func TestSetupWritesStoragePathConfiguration(t *testing.T) {
	home := givenUserHome(t)
	location := filepath.Join(home, "managed-tools")

	if _, err := environment.Setup(location); err != nil {
		t.Fatalf("Setup() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".mdev", "config.yaml"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	expected := "storage_path: " + location + "\n"
	if string(data) != expected {
		t.Fatalf("config contents = %q, want %q", data, expected)
	}
}

func TestConfigLoadReadsStoragePath(t *testing.T) {
	home := givenUserHome(t)
	configDir := filepath.Join(home, ".mdev")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	location := filepath.Join(home, "mdev")
	data := []byte("storage_path: " + location + "\n")
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.StoragePath != location {
		t.Fatalf("StoragePath = %q, want %q", cfg.StoragePath, location)
	}
}

func TestConfigLoadRejectsLegacyExternalDriveKey(t *testing.T) {
	home := givenUserHome(t)
	configDir := filepath.Join(home, ".mdev")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(configDir, "config.yaml"),
		[]byte("external_drive: /legacy/path\n"),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	if _, err := config.Load(); !errors.Is(err, config.ErrStoragePathRequired) {
		t.Fatalf("Load() error = %v, want ErrStoragePathRequired", err)
	}
}

func TestToolStorageIsDirectlyBelowStoragePath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "mdev")
	env := environment.New(root)

	got := storage.ToolDir(env, "maven")
	want := filepath.Join(root, "maven")
	if got != want {
		t.Fatalf("ToolDir() = %q, want %q", got, want)
	}
	if got == filepath.Join(root, "data", "maven") {
		t.Fatal("ToolDir() retained obsolete data directory")
	}
}

func TestValidateStoragePathExpandsHomeDirectory(t *testing.T) {
	home := givenUserHome(t)

	got, err := environment.ValidateStoragePath("~/Documents/mdev")
	if err != nil {
		t.Fatalf("ValidateStoragePath() error = %v", err)
	}

	want := filepath.Join(home, "Documents", "mdev")
	if got != want {
		t.Fatalf("ValidateStoragePath() = %q, want %q", got, want)
	}
}

func TestValidateStoragePathRejectsBroadLocations(t *testing.T) {
	home := givenUserHome(t)

	for _, path := range []string{"", string(filepath.Separator), home} {
		if _, err := environment.ValidateStoragePath(path); err == nil {
			t.Fatalf("ValidateStoragePath(%q) unexpectedly succeeded", path)
		}
	}
}
