package podmandesktop

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

func TestMetadataKeepsDesktopSeparateFromPodman(t *testing.T) {
	tool := &PodmanDesktop{}
	if got := tool.Dependencies(); len(got) != 1 || got[0] != "podman" {
		t.Fatalf("dependencies = %v, want [podman]", got)
	}
	if got := tool.StorageDir(environment.New(t.TempDir())); got != "" {
		t.Fatalf("storage directory = %q, want empty", got)
	}
}

func TestInstallationStatusDetectsExistingUserApplication(t *testing.T) {
	home := t.TempDir()
	app := filepath.Join(home, "Applications", "Podman Desktop.app")
	if err := os.MkdirAll(app, 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("PATH", t.TempDir())
	if !(&PodmanDesktop{}).InstallationStatus(environment.New(t.TempDir())) {
		t.Fatal("existing Desktop application was not detected")
	}
}

func TestInstallUsesDesktopCask(t *testing.T) {
	bin, logPath := fakeBrew(t, "")
	t.Setenv("PATH", bin)
	t.Setenv("MDEV_TEST_BREW_LOG", logPath)
	if err := (&PodmanDesktop{}).InstallContext(context.Background(), environment.New(t.TempDir())); err != nil {
		t.Fatal(err)
	}
	assertBrewLog(t, logPath, "install --cask podman-desktop")
}

func TestUninstallPreservesApplicationNotOwnedByHomebrew(t *testing.T) {
	bin, logPath := fakeBrew(t, "exit 1")
	t.Setenv("PATH", bin)
	t.Setenv("MDEV_TEST_BREW_LOG", logPath)
	if err := (&PodmanDesktop{}).Uninstall(environment.New(t.TempDir())); err != nil {
		t.Fatal(err)
	}
	assertBrewLog(t, logPath, "list --cask podman-desktop")
}

func fakeBrew(t *testing.T, listResult string) (string, string) {
	t.Helper()
	if listResult == "" {
		listResult = ":"
	}
	bin := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "brew.log")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$MDEV_TEST_BREW_LOG\"\nif [ \"$1\" = list ]; then " + listResult + "; fi\n"
	if err := os.WriteFile(filepath.Join(bin, "brew"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	return bin, logPath
}

func assertBrewLog(t *testing.T, path, want string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(content)); got != want {
		t.Fatalf("brew calls = %q, want %q", got, want)
	}
}
