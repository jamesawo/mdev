package podmancompose

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

func TestMetadataKeepsComposeSeparateFromPodman(t *testing.T) {
	tool := &PodmanCompose{}
	if got := tool.Dependencies(); len(got) != 1 || got[0] != "podman" {
		t.Fatalf("dependencies = %v, want [podman]", got)
	}
	if got := tool.StorageDir(environment.New(t.TempDir())); got != "" {
		t.Fatalf("storage directory = %q, want empty", got)
	}
}

func TestInstallationStatusRequiresHomebrewOwnershipAndWorkingProvider(t *testing.T) {
	bin, logPath := fakeCommands(t, true)
	t.Setenv("PATH", bin)
	t.Setenv("MDEV_TEST_BREW_LOG", logPath)
	if !(&PodmanCompose{}).InstallationStatus(environment.New(t.TempDir())) {
		t.Fatal("Homebrew-managed working provider was not detected")
	}

	bin, logPath = fakeCommands(t, false)
	t.Setenv("PATH", bin)
	t.Setenv("MDEV_TEST_BREW_LOG", logPath)
	if (&PodmanCompose{}).InstallationStatus(environment.New(t.TempDir())) {
		t.Fatal("unmanaged provider counted as mdev-managed completion")
	}
}

func TestInstallUsesHomebrewFormula(t *testing.T) {
	bin, logPath := fakeCommands(t, true)
	t.Setenv("PATH", bin)
	t.Setenv("MDEV_TEST_BREW_LOG", logPath)
	if err := (&PodmanCompose{}).InstallContext(context.Background(), environment.New(t.TempDir())); err != nil {
		t.Fatal(err)
	}
	assertLog(t, logPath, "install podman-compose")
}

func TestUninstallRemovesOnlyOwnedFormula(t *testing.T) {
	bin, logPath := fakeCommands(t, true)
	t.Setenv("PATH", bin)
	t.Setenv("MDEV_TEST_BREW_LOG", logPath)
	if err := (&PodmanCompose{}).Uninstall(environment.New(t.TempDir())); err != nil {
		t.Fatal(err)
	}
	assertLog(t, logPath, "list podman-compose\nuninstall podman-compose")
}

func fakeCommands(t *testing.T, owned bool) (string, string) {
	t.Helper()
	bin := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "brew.log")
	listResult := "exit 1"
	if owned {
		listResult = "exit 0"
	}
	brewScript := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$MDEV_TEST_BREW_LOG\"\nif [ \"$1\" = list ]; then " + listResult + "; fi\n"
	if err := os.WriteFile(filepath.Join(bin, "brew"), []byte(brewScript), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bin, formulaName), []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	return bin, logPath
}

func assertLog(t *testing.T, path, want string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(content)); got != want {
		t.Fatalf("brew calls = %q, want %q", got, want)
	}
}
