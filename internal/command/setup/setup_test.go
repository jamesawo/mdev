package setup

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

func TestInterruptIsCancellation(t *testing.T) {
	if err := interruptionError(terminal.InterruptErr); !errors.Is(err, environment.ErrSetupCancelled) {
		t.Fatalf("interruptionError() = %v", err)
	}
}

func TestInteractiveAlreadyConfiguredEndsWithListRecommendation(t *testing.T) {
	env := environment.New("/storage/mdev")
	deps := stubDependencies()
	deps.setupInteractive = func() (*environment.Environment, error) {
		return env, environment.ErrAlreadyConfigured
	}
	var out recordingOutput

	if err := run(Options{}, deps, &out); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"section:" + messages.SetupTitle,
		"info:" + messages.SetupComplete,
		"info:" + messages.SetupStorage("display:/storage/mdev"),
		"blank",
		"info:" + messages.SetupNextStep,
		"command:" + messages.SetupListCommand,
	}
	if !reflect.DeepEqual(out.calls, want) {
		t.Fatalf("output calls = %#v, want %#v", out.calls, want)
	}
	if got := out.calls[len(out.calls)-1]; got != "command:mdev list" {
		t.Fatalf("last output = %q", got)
	}
}

func TestInteractiveCancellationIsSuccessful(t *testing.T) {
	deps := stubDependencies()
	deps.setupInteractive = func() (*environment.Environment, error) {
		return nil, environment.ErrSetupCancelled
	}
	var out recordingOutput
	if err := run(Options{}, deps, &out); err != nil {
		t.Fatal(err)
	}
	want := []string{"section:" + messages.SetupTitle, "info:" + messages.SetupCancelled}
	if !reflect.DeepEqual(out.calls, want) {
		t.Fatalf("output calls = %#v, want %#v", out.calls, want)
	}
}

func TestInteractiveFailurePreservesCause(t *testing.T) {
	wantErr := errors.New("interactive failure")
	deps := stubDependencies()
	deps.setupInteractive = func() (*environment.Environment, error) { return nil, wantErr }
	if err := run(Options{}, deps, &recordingOutput{}); !errors.Is(err, wantErr) || !strings.Contains(err.Error(), messages.SetupFailed) {
		t.Fatalf("run() error = %v", err)
	}
}

func TestNonInteractiveSetupCreatesCanonicalStorageAndConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	parent := filepath.Join(t.TempDir(), "parent with spaces")
	var out recordingOutput
	if err := run(Options{StoragePath: parent}, productionDependencies(), &out); err != nil {
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
	if got := out.calls[len(out.calls)-1]; got != "command:mdev list" {
		t.Fatalf("last output = %q", got)
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
	replacement := filepath.Join(t.TempDir(), "replacement")
	err := run(Options{StoragePath: replacement}, productionDependencies(), &recordingOutput{})
	if err == nil || !strings.Contains(err.Error(), existing) {
		t.Fatalf("run() error = %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(replacement, "mdev")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("replacement storage was created: %v", statErr)
	}
}

func TestNonInteractiveResolutionFailurePreservesCause(t *testing.T) {
	wantErr := errors.New("resolve failure")
	deps := stubDependencies()
	deps.resolveStorage = func(string) (string, string, error) { return "", "", wantErr }
	if err := run(Options{StoragePath: "somewhere"}, deps, &recordingOutput{}); !errors.Is(err, wantErr) {
		t.Fatalf("run() error = %v", err)
	}
}

func stubDependencies() dependencies {
	return dependencies{
		setupInteractive: func() (*environment.Environment, error) {
			return environment.New("/storage/mdev"), nil
		},
		existing:       func() (*environment.Environment, bool, error) { return nil, false, nil },
		resolveStorage: func(path string) (string, string, error) { return path, "", nil },
		setupResolved: func(path string) (*environment.Environment, error) {
			return environment.New(path), nil
		},
		displayStoragePath: func(path string) string { return "display:" + path },
	}
}

type recordingOutput struct {
	calls []string
}

func (o *recordingOutput) Section(text string) { o.calls = append(o.calls, "section:"+text) }
func (o *recordingOutput) Info(text string)    { o.calls = append(o.calls, "info:"+text) }
func (o *recordingOutput) Blank()              { o.calls = append(o.calls, "blank") }
func (o *recordingOutput) Command(text string) { o.calls = append(o.calls, "command:"+text) }
