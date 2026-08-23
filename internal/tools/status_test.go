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

func TestManagedSymlinkStatusRequiresExpectedTarget(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	other := filepath.Join(root, "other")
	if err := os.Mkdir(target, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(other, 0755); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(root, "source")
	if err := os.Symlink(other, source); err != nil {
		t.Fatal(err)
	}
	if installed, err := ManagedSymlinkStatus(source, target); err != nil || installed {
		t.Fatalf("wrong symlink status = %v, %v", installed, err)
	}
	if err := os.Remove(source); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, source); err != nil {
		t.Fatal(err)
	}
	if installed, err := ManagedSymlinkStatus(source, target); err != nil || !installed {
		t.Fatalf("managed symlink status = %v, %v", installed, err)
	}
}

func TestManagedSymlinkStatusTreatsExpectedDanglingLinkAsIncomplete(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "missing-target")
	source := filepath.Join(root, "source")
	if err := os.Symlink(target, source); err != nil {
		t.Fatal(err)
	}
	installed, err := ManagedSymlinkStatus(source, target)
	if err != nil {
		t.Fatalf("status returned an error for recoverable state: %v", err)
	}
	if installed {
		t.Fatal("dangling managed symlink was reported installed")
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
