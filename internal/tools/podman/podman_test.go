package podman

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

func TestInstallProvisionsCLIAndDesktopInOrder(t *testing.T) {
	bin := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "brew.log")
	brewPath := filepath.Join(bin, "brew")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$MDEV_TEST_BREW_LOG\"\n"
	if err := os.WriteFile(brewPath, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	t.Setenv("MDEV_TEST_BREW_LOG", logPath)

	if err := (&Podman{}).InstallContext(context.Background(), environment.New(t.TempDir())); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.FieldsFunc(strings.TrimSpace(string(output)), func(r rune) bool { return r == '\n' })
	want := []string{"install podman", "install --cask podman-desktop"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("brew calls = %q, want %q", got, want)
	}
}

func TestConfigureRelocatesMachineStateBeforeInitialization(t *testing.T) {
	home := t.TempDir()
	bin := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "podman.log")
	podmanPath := filepath.Join(bin, "podman")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$MDEV_TEST_PODMAN_LOG\"\n"
	if err := os.WriteFile(podmanPath, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	legacyState := podmanMachineDir(home)
	if err := os.MkdirAll(legacyState, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(legacyState, "existing"), []byte("state"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("PATH", bin)
	t.Setenv("MDEV_TEST_PODMAN_LOG", logPath)
	env := environment.New(filepath.Join(t.TempDir(), "mdev"))

	if err := (&Podman{}).ConfigureContext(context.Background(), env); err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(legacyState)
	if err != nil {
		t.Fatal(err)
	}
	wantTarget := filepath.Join(env.StoragePath, "podman")
	wantResolved, err := filepath.EvalSymlinks(wantTarget)
	if err != nil {
		t.Fatal(err)
	}
	if resolved != wantResolved {
		t.Fatalf("machine state resolves to %q, want %q", resolved, wantResolved)
	}
	if _, err := os.Stat(filepath.Join(wantTarget, "existing")); err != nil {
		t.Fatalf("existing machine state was not preserved: %v", err)
	}
	output, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(output)); got != "machine init" {
		t.Fatalf("podman call = %q", got)
	}
}
