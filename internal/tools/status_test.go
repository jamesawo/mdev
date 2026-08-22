package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

func TestCommandInstallationStatus(t *testing.T) {
	dir := t.TempDir()
	installed := filepath.Join(dir, "installed")
	failed := filepath.Join(dir, "failed")
	if err := os.WriteFile(installed, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(failed, []byte("#!/bin/sh\necho broken >&2\nexit 2\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	if got, err := CommandInstallationStatus("installed"); err != nil || !got {
		t.Fatalf("installed status = %v, %v", got, err)
	}
	if got, err := CommandInstallationStatus("missing"); err != nil || got {
		t.Fatalf("missing status = %v, %v", got, err)
	}
	got, err := CommandInstallationStatus("failed")
	if err == nil || got || !strings.Contains(err.Error(), "broken") {
		t.Fatalf("failed status = %v, %v", got, err)
	}
}

func TestInstallationStatusFallsBackToToolContract(t *testing.T) {
	tool := fallbackStatusTool{installed: true}
	installed, err := InstallationStatus(tool, nil)
	if err != nil || !installed {
		t.Fatalf("InstallationStatus() = %v, %v", installed, err)
	}
}

type fallbackStatusTool struct {
	installed bool
}

func (t fallbackStatusTool) Name() string                               { return "fallback" }
func (t fallbackStatusTool) Description() string                        { return "" }
func (t fallbackStatusTool) Dependencies() []string                     { return nil }
func (t fallbackStatusTool) IsInstalled(*environment.Environment) bool  { return t.installed }
func (t fallbackStatusTool) Install(*environment.Environment) error     { return nil }
func (t fallbackStatusTool) Configure(*environment.Environment) error   { return nil }
func (t fallbackStatusTool) Verify(*environment.Environment) error      { return nil }
func (t fallbackStatusTool) Uninstall(*environment.Environment) error   { return nil }
func (t fallbackStatusTool) StorageDir(*environment.Environment) string { return "" }
