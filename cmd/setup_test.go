package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
)

func TestSetupHelpOmitsYesAndDocumentsStoragePath(t *testing.T) {
	var output bytes.Buffer
	setupCmd.SetOut(&output)
	t.Cleanup(func() { setupCmd.SetOut(nil) })
	if err := setupCmd.Help(); err != nil {
		t.Fatal(err)
	}
	help := output.String()
	if strings.Contains(help, "--yes") {
		t.Fatalf("setup help includes unsupported --yes flag:\n%s", help)
	}
	if !strings.Contains(help, "--storage-path") {
		t.Fatalf("setup help omits --storage-path:\n%s", help)
	}
}

func TestSetupRejectsInheritedYes(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("yes")
	if err := flag.Value.Set("true"); err != nil {
		t.Fatal(err)
	}
	flag.Changed = true
	t.Cleanup(func() {
		_ = flag.Value.Set("false")
		flag.Changed = false
	})
	if err := setupCmd.PreRunE(setupCmd, nil); err == nil || !strings.Contains(err.Error(), "does not support") {
		t.Fatalf("PreRunE() error = %v", err)
	}
}

func TestNonInteractiveSetupCreatesCanonicalStorageAndConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	parent := filepath.Join(t.TempDir(), "parent with spaces")
	setupStoragePath = parent
	t.Cleanup(func() { setupStoragePath = "" })
	if err := runSetup(setupCmd, nil); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(cfg.StoragePath) != "mdev" {
		t.Fatalf("storage path = %q", cfg.StoragePath)
	}
	if info, err := os.Stat(cfg.StoragePath); err != nil || !info.IsDir() {
		t.Fatalf("storage path was not created: %v", err)
	}
}

func TestNonInteractiveSetupRefusesExistingConfiguration(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	existing := filepath.Join(t.TempDir(), "existing", "mdev")
	if err := os.MkdirAll(existing, 0755); err != nil {
		t.Fatal(err)
	}
	if err := config.Save(config.Config{StoragePath: existing}); err != nil {
		t.Fatal(err)
	}
	setupStoragePath = filepath.Join(t.TempDir(), "replacement")
	t.Cleanup(func() { setupStoragePath = "" })
	err := runSetup(setupCmd, nil)
	if err == nil || !strings.Contains(err.Error(), existing) {
		t.Fatalf("runSetup() error = %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(setupStoragePath, "mdev")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("replacement storage was created: %v", statErr)
	}
}
