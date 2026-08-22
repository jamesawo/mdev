package environment

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jamesawo/mdev/internal/infrastructure/config"
)

func TestResolveStoragePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	working := t.TempDir()
	oldWorking, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(working); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWorking) })

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "home expansion", input: "~/Documents", want: filepath.Join(home, "Documents", "mdev")},
		{name: "relative path", input: "work area", want: filepath.Join(working, "work area", "mdev")},
		{name: "clean path", input: filepath.Join(working, "one", "..", "two"), want: filepath.Join(working, "two", "mdev")},
		{name: "already suffixed", input: filepath.Join(working, "mdev"), want: filepath.Join(working, "mdev")},
		{name: "no variable expansion", input: "$HOME/place", want: filepath.Join(working, "$HOME", "place", "mdev")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := ResolveStoragePath(test.input)
			if err != nil {
				t.Fatal(err)
			}
			want := canonical(t, test.want)
			if got != want {
				t.Fatalf("ResolveStoragePath(%q) = %q, want %q", test.input, got, want)
			}
		})
	}
}

func TestResolveStoragePathCanonicalizesSymlink(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	target := t.TempDir()
	link := filepath.Join(home, "storage")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	got, source, err := ResolveStoragePath(link)
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Join(canonical(t, target), "mdev") {
		t.Fatalf("resolved path = %q", got)
	}
	if source != link {
		t.Fatalf("symlink source = %q, want %q", source, link)
	}
}

func TestDisplayPathShortensCanonicalHomePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	canonicalHome := canonical(t, home)
	if got := DisplayPath(filepath.Join(canonicalHome, "mdev")); got != filepath.Join("~", "mdev") {
		t.Fatalf("DisplayPath() = %q", got)
	}
}

func TestResolveStoragePathRejectsFileAndBroadPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	file := filepath.Join(t.TempDir(), "mdev")
	if err := os.WriteFile(file, []byte("occupied"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"", "/", home, file} {
		if _, _, err := ResolveStoragePath(path); err == nil {
			t.Fatalf("ResolveStoragePath(%q) unexpectedly succeeded", path)
		}
	}
}

func TestSetupRollsBackStorageWhenPersistenceFails(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	parent := filepath.Join(t.TempDir(), "new-parent")
	storage := filepath.Join(parent, "mdev")
	originalSave := saveConfig
	saveConfig = func(config.Config) error { return errors.New("injected persistence failure") }
	t.Cleanup(func() { saveConfig = originalSave })

	if _, err := SetupResolved(storage); err == nil || !strings.Contains(err.Error(), "persistence failure") {
		t.Fatalf("SetupResolved() error = %v", err)
	}
	if _, err := os.Stat(parent); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("created directories were not rolled back: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, ".mdev")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("configuration artifacts remain: %v", err)
	}
}

func TestCreateStorageRootOwnsEveryNewDirectory(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "new", "nested")
	storage := filepath.Join(parent, "mdev")
	var owned []string
	originalOwnPath := ownPath
	ownPath = func(path string) error {
		owned = append(owned, path)
		return nil
	}
	t.Cleanup(func() { ownPath = originalOwnPath })

	created, err := createStorageRoot(storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(owned) != len(created) {
		t.Fatalf("owned paths = %#v, created paths = %#v", owned, created)
	}
	for index := range created {
		if owned[index] != created[index] {
			t.Fatalf("owned paths = %#v, created paths = %#v", owned, created)
		}
	}
}

func TestSetupResolvedRejectsUnresolvedOrBroadPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	for _, path := range []string{"relative/mdev", "/", filepath.Join(home, "other")} {
		if _, err := SetupResolved(path); err == nil {
			t.Fatalf("SetupResolved(%q) unexpectedly succeeded", path)
		}
	}
}

func TestSetupNeverReplacesExistingConfiguration(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	first := filepath.Join(t.TempDir(), "first", "mdev")
	second := filepath.Join(t.TempDir(), "second", "mdev")
	if _, err := SetupResolved(first); err != nil {
		t.Fatal(err)
	}
	if _, err := SetupResolved(second); !errors.Is(err, ErrAlreadyConfigured) {
		t.Fatalf("second setup error = %v", err)
	}
	if _, err := os.Stat(second); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("second storage path was created: %v", err)
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.StoragePath != canonical(t, first) {
		t.Fatalf("configuration changed to %q", cfg.StoragePath)
	}
}

func TestExistingLeavesMalformedConfigurationUntouched(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	dir := filepath.Join(home, ".mdev")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.yaml")
	original := []byte("storage_path: [broken")
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Existing(); err == nil || !strings.Contains(err.Error(), "fix or remove") {
		t.Fatalf("Existing() error = %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("malformed configuration changed to %q", got)
	}
}

func TestInterruptIsCancellation(t *testing.T) {
	if err := interruptionError(terminal.InterruptErr); !errors.Is(err, ErrSetupCancelled) {
		t.Fatalf("interruptionError() = %v", err)
	}
}

func TestPermissionErrorSuggestsAnotherLocationOrSudo(t *testing.T) {
	err := actionablePathError("/protected/mdev", os.ErrPermission)
	if !strings.Contains(err.Error(), "choose another location") || !strings.Contains(err.Error(), "sudo") {
		t.Fatalf("actionablePathError() = %q", err)
	}
}

func TestDefaultLocationUsesMdevDirectoryInUserHome(t *testing.T) {
	home := givenEnvironmentHome(t)

	location, err := DefaultLocation()
	if err != nil {
		t.Fatalf("expected default location to resolve: %v", err)
	}

	expected := filepath.Join(home, "mdev")
	if location != expected {
		t.Fatalf("expected default location %q, got %q", expected, location)
	}
}

func TestSetupPreservesExistingContents(t *testing.T) {
	home := givenEnvironmentHome(t)
	location := filepath.Join(home, "mdev")
	existingFile := filepath.Join(location, "existing.txt")

	if err := os.MkdirAll(location, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(existingFile, []byte("keep me"), 0644); err != nil {
		t.Fatal(err)
	}

	env, err := Setup(location)
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
	if cfg.StoragePath != env.StoragePath {
		t.Fatalf("configured location = %q, want %q", cfg.StoragePath, env.StoragePath)
	}
}

func TestSetupWritesStoragePathConfiguration(t *testing.T) {
	home := givenEnvironmentHome(t)
	location := filepath.Join(home, "managed-tools")

	if _, err := Setup(location); err != nil {
		t.Fatalf("Setup() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".mdev", "config.yaml"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	resolved, err := ValidateStoragePath(location)
	if err != nil {
		t.Fatal(err)
	}
	expected := "storage_path: " + resolved + "\n"
	if string(data) != expected {
		t.Fatalf("config contents = %q, want %q", data, expected)
	}
}

func TestValidateStoragePathExpandsHomeDirectory(t *testing.T) {
	home := givenEnvironmentHome(t)

	got, err := ValidateStoragePath("~/Documents/mdev")
	if err != nil {
		t.Fatalf("ValidateStoragePath() error = %v", err)
	}

	want, err := filepath.EvalSymlinks(home)
	if err != nil {
		t.Fatal(err)
	}
	want = filepath.Join(want, "Documents", "mdev")
	if got != want {
		t.Fatalf("ValidateStoragePath() = %q, want %q", got, want)
	}
}

func TestValidateStoragePathRejectsBroadLocations(t *testing.T) {
	home := givenEnvironmentHome(t)

	for _, path := range []string{"", string(filepath.Separator), home} {
		if _, err := ValidateStoragePath(path); err == nil {
			t.Fatalf("ValidateStoragePath(%q) unexpectedly succeeded", path)
		}
	}
}

func givenEnvironmentHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	return home
}

func canonical(t *testing.T, path string) string {
	t.Helper()
	ancestor := path
	var suffix []string
	for {
		if _, err := os.Stat(ancestor); err == nil {
			break
		}
		suffix = append(suffix, filepath.Base(ancestor))
		ancestor = filepath.Dir(ancestor)
	}
	resolved, err := filepath.EvalSymlinks(ancestor)
	if err != nil {
		t.Fatal(err)
	}
	for i := len(suffix) - 1; i >= 0; i-- {
		resolved = filepath.Join(resolved, suffix[i])
	}
	return resolved
}
