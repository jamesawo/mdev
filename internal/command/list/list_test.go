package list

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
	"github.com/jamesawo/mdev/internal/tools"
	_ "github.com/jamesawo/mdev/internal/tools/gradle"
	_ "github.com/jamesawo/mdev/internal/tools/java"
	_ "github.com/jamesawo/mdev/internal/tools/maven"
	_ "github.com/jamesawo/mdev/internal/tools/nvm"
	_ "github.com/jamesawo/mdev/internal/tools/podman"
	_ "github.com/jamesawo/mdev/internal/tools/sdkman"
)

func TestProductionDependenciesIncludeEveryRegisteredEntry(t *testing.T) {
	deps := productionDependencies()
	env := environment.New(t.TempDir())

	gotSystems := checkNames(deps.systemChecks())
	wantSystems := registeredPrerequisiteNames(prerequisites.List())
	if !reflect.DeepEqual(gotSystems, wantSystems) {
		t.Fatalf("system checks = %#v, want registered prerequisites %#v", gotSystems, wantSystems)
	}
	gotTools := checkNames(deps.toolChecks(env))
	wantTools := registeredToolNames(tools.List())
	if !reflect.DeepEqual(gotTools, wantTools) {
		t.Fatalf("tool checks = %#v, want registered tools %#v", gotTools, wantTools)
	}
}

func TestRunGroupsSortsAndReportsEveryStatus(t *testing.T) {
	storage := t.TempDir()
	var checked []string
	checkFor := func(name string, installed bool, err error) check {
		return check{name: name, verify: func() (bool, error) {
			checked = append(checked, name)
			return installed, err
		}}
	}
	deps := validDependencies(storage)
	deps.systemChecks = func() []check {
		return []check{
			checkFor("zeta", false, nil),
			checkFor("Alpha", true, nil),
		}
	}
	verificationErr := errors.New("permission denied")
	deps.toolChecks = func(*environment.Environment) []check {
		return []check{
			checkFor("gamma", true, nil),
			checkFor("beta", false, verificationErr),
		}
	}

	var output bytes.Buffer
	err := run(&output, deps)
	var unknownErr *UnknownStatusError
	if !errors.As(err, &unknownErr) {
		t.Fatalf("run() error = %v, want UnknownStatusError", err)
	}
	if !errors.Is(err, verificationErr) {
		t.Fatalf("run() error does not preserve verification failure: %v", err)
	}
	wantChecks := []string{"Alpha", "zeta", "beta", "gamma"}
	if !reflect.DeepEqual(checked, wantChecks) {
		t.Fatalf("checks = %#v, want %#v", checked, wantChecks)
	}
	want := "system tools\n" +
		"  ✓ Alpha  installed\n" +
		"  ○ zeta   missing\n" +
		"\n" +
		"tools\n" +
		"  ? beta   unknown\n" +
		"  ✓ gamma  installed\n" +
		"\n" +
		"could not determine beta status: permission denied\n"
	if output.String() != want {
		t.Fatalf("output:\n%s\nwant:\n%s", output.String(), want)
	}
}

func TestRunOmitsEmptySections(t *testing.T) {
	deps := validDependencies(t.TempDir())
	deps.systemChecks = func() []check { return nil }
	deps.toolChecks = func(*environment.Environment) []check {
		return []check{{name: "go", verify: func() (bool, error) { return true, nil }}}
	}

	var output bytes.Buffer
	if err := run(&output, deps); err != nil {
		t.Fatal(err)
	}
	if output.String() != "tools\n  ✓ go  installed\n" {
		t.Fatalf("output = %q", output.String())
	}
}

func TestRunRequiresConfigurationWithoutChangingHome(t *testing.T) {
	home := t.TempDir()
	sentinel := filepath.Join(home, "sentinel")
	if err := os.WriteFile(sentinel, []byte("unchanged"), 0600); err != nil {
		t.Fatal(err)
	}
	deps := dependencies{
		loadEnvironment: func() (*environment.Environment, bool, error) { return nil, false, nil },
		stat: func(string) (os.FileInfo, error) {
			t.Fatal("storage was checked without configuration")
			return nil, nil
		},
	}

	var output bytes.Buffer
	err := run(&output, deps)
	if err == nil || !strings.Contains(err.Error(), "mdev setup") {
		t.Fatalf("run() error = %v", err)
	}
	if output.Len() != 0 {
		t.Fatalf("output = %q", output.String())
	}
	data, readErr := os.ReadFile(sentinel)
	if readErr != nil || string(data) != "unchanged" {
		t.Fatalf("sentinel changed: data=%q error=%v", data, readErr)
	}
	if _, statErr := os.Stat(filepath.Join(home, ".mdev")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("configuration was created: %v", statErr)
	}
}

func TestRunPreservesConfigurationError(t *testing.T) {
	configErr := errors.New("malformed configuration")
	deps := dependencies{
		loadEnvironment: func() (*environment.Environment, bool, error) { return nil, false, configErr },
	}
	if err := run(&bytes.Buffer{}, deps); !errors.Is(err, configErr) {
		t.Fatalf("run() error = %v, want %v", err, configErr)
	}
}

func TestRunLeavesMalformedConfigurationUntouched(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	configDir := filepath.Join(home, ".mdev")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	original := []byte("storage_path: [broken")
	if err := os.WriteFile(configPath, original, 0600); err != nil {
		t.Fatal(err)
	}
	deps := productionDependencies()
	deps.systemChecks = func() []check { return nil }
	deps.toolChecks = func(*environment.Environment) []check { return nil }

	err := run(&bytes.Buffer{}, deps)
	if err == nil || !strings.Contains(err.Error(), "configuration cannot be read") {
		t.Fatalf("run() error = %v", err)
	}
	got, readErr := os.ReadFile(configPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(got, original) {
		t.Fatalf("configuration changed to %q", got)
	}
}

func TestRunLeavesValidConfigurationAndStorageUntouched(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	storage := t.TempDir()
	configDir := filepath.Join(home, ".mdev")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	configData := []byte("storage_path: " + storage + "\n")
	if err := os.WriteFile(configPath, configData, 0600); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(storage, "sentinel")
	if err := os.WriteFile(sentinel, []byte("unchanged"), 0600); err != nil {
		t.Fatal(err)
	}
	deps := productionDependencies()
	deps.systemChecks = func() []check { return nil }
	deps.toolChecks = func(*environment.Environment) []check { return nil }

	if err := run(&bytes.Buffer{}, deps); err != nil {
		t.Fatal(err)
	}
	gotConfig, err := os.ReadFile(configPath)
	if err != nil || !bytes.Equal(gotConfig, configData) {
		t.Fatalf("configuration changed: data=%q error=%v", gotConfig, err)
	}
	gotSentinel, err := os.ReadFile(sentinel)
	if err != nil || string(gotSentinel) != "unchanged" {
		t.Fatalf("storage changed: data=%q error=%v", gotSentinel, err)
	}
}

func TestRunRejectsUnavailableStorageBeforeChecks(t *testing.T) {
	storage := filepath.Join(t.TempDir(), "disconnected")
	checked := false
	deps := validDependencies(storage)
	deps.systemChecks = func() []check {
		checked = true
		return nil
	}

	err := run(&bytes.Buffer{}, deps)
	if err == nil || !strings.Contains(err.Error(), storage) || !strings.Contains(err.Error(), "unavailable") {
		t.Fatalf("run() error = %v", err)
	}
	if checked {
		t.Fatal("tool checks ran with unavailable storage")
	}
	if _, statErr := os.Stat(storage); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("storage was created: %v", statErr)
	}
}

func TestRunRejectsStorageFile(t *testing.T) {
	storage := filepath.Join(t.TempDir(), "storage")
	if err := os.WriteFile(storage, nil, 0600); err != nil {
		t.Fatal(err)
	}
	deps := validDependencies(storage)
	err := run(&bytes.Buffer{}, deps)
	if err == nil || !strings.Contains(err.Error(), "expected a directory") {
		t.Fatalf("run() error = %v", err)
	}
}

func validDependencies(storage string) dependencies {
	return dependencies{
		loadEnvironment: func() (*environment.Environment, bool, error) {
			return environment.New(storage), true, nil
		},
		stat:         os.Stat,
		systemChecks: func() []check { return nil },
		toolChecks:   func(*environment.Environment) []check { return nil },
	}
}

func checkNames(checks []check) []string {
	names := make([]string, 0, len(checks))
	for _, item := range checks {
		names = append(names, item.name)
	}
	sort.Strings(names)
	return names
}

func registeredPrerequisiteNames(registered []prerequisites.Prerequisite) []string {
	names := make([]string, 0, len(registered))
	for _, item := range registered {
		names = append(names, item.Name())
	}
	sort.Strings(names)
	return names
}

func registeredToolNames(registered []tools.Tool) []string {
	names := make([]string, 0, len(registered))
	for _, item := range registered {
		names = append(names, item.Name())
	}
	sort.Strings(names)
	return names
}
