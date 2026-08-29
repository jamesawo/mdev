package uninstall

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/command"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
)

type fakeTool struct {
	name    string
	storage string
}

func (t fakeTool) Name() string                              { return t.name }
func (t fakeTool) Description() string                       { return t.name }
func (t fakeTool) Dependencies() []string                    { return nil }
func (t fakeTool) IsInstalled(*environment.Environment) bool { return true }
func (t fakeTool) Install(*environment.Environment) error    { return nil }
func (t fakeTool) Configure(*environment.Environment) error  { return nil }
func (t fakeTool) Verify(*environment.Environment) error     { return nil }
func (t fakeTool) Uninstall(*environment.Environment) error  { return nil }
func (t fakeTool) StorageDir(*environment.Environment) string {
	return t.storage
}

func testWorkflowDependencies(t *testing.T, plan []string, registered map[string]tools.Tool, installed map[string]bool) workflowDependencies {
	t.Helper()
	return workflowDependencies{
		loadEnvironment: func() (*environment.Environment, error) { return environment.New(t.TempDir()), nil },
		buildPlan:       func(string) ([]string, error) { return plan, nil },
		status:          func(name string, _ *environment.Environment) (bool, error) { return installed[name], nil },
		getTool:         func(name string) (tools.Tool, bool) { tool, ok := registered[name]; return tool, ok },
		stat:            os.Stat,
		removeAll:       os.RemoveAll,
		uninstall: func(_ context.Context, tool tools.Tool, _ *environment.Environment) error {
			installed[tool.Name()] = false
			return nil
		},
		newReporter: func(out io.Writer) progressReporter { return newTextProgressReporter(out) },
	}
}

func TestRunPresentsDependencyRemovalAndStorageCleanupInOrder(t *testing.T) {
	storage := filepath.Join(t.TempDir(), "podman")
	if err := os.MkdirAll(storage, 0755); err != nil {
		t.Fatal(err)
	}
	registered := map[string]tools.Tool{
		"podman-compose": fakeTool{name: "podman-compose"},
		"podman":         fakeTool{name: "podman", storage: storage},
	}
	installed := map[string]bool{"podman-compose": true, "podman": true}
	deps := testWorkflowDependencies(t, []string{"podman-compose", "podman"}, registered, installed)
	var output bytes.Buffer
	if err := run(context.Background(), Streams{In: strings.NewReader(""), Out: &output}, Options{Tool: "podman", AssumeYes: true}, deps); err != nil {
		t.Fatal(err)
	}
	want := "dependency warning\n" +
		"  podman is required by:\n" +
		"    podman-compose\n" +
		"uninstall plan\n" +
		"  podman-compose\n" +
		"  podman\n" +
		"directories to be removed\n" +
		"  " + storage + "\n" +
		"uninstalling podman-compose... ✓\n" +
		"uninstalling podman... ✓\n" +
		"cleaning podman storage... ✓\n" +
		"podman-compose removed.\n" +
		"podman removed.\n"
	if output.String() != want {
		t.Fatalf("output:\n%s\nwant:\n%s", output.String(), want)
	}
	if _, err := os.Stat(storage); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("managed storage still exists: %v", err)
	}
}

func TestRunCapturesProviderOutputAndRetainsFailureDiagnostic(t *testing.T) {
	tool := fakeTool{name: "example"}
	installed := map[string]bool{"example": true}
	deps := testWorkflowDependencies(t, []string{"example"}, map[string]tools.Tool{"example": tool}, installed)
	deps.uninstall = func(ctx context.Context, _ tools.Tool, _ *environment.Environment) error {
		return command.RunContext(ctx, "sh", "-c", "printf 'provider success chatter'; printf 'provider removal failed' >&2; exit 6")
	}
	var output bytes.Buffer
	err := run(context.Background(), Streams{In: strings.NewReader(""), Out: &output}, Options{Tool: "example", AssumeYes: true}, deps)
	if err == nil || !strings.Contains(err.Error(), "uninstall example") || !strings.Contains(err.Error(), "provider removal failed") {
		t.Fatalf("error = %v, want tool and provider diagnostic", err)
	}
	if strings.Contains(output.String(), "provider") {
		t.Fatalf("provider output leaked into presentation: %q", output.String())
	}
	if !strings.Contains(output.String(), "uninstalling example... ✗") || strings.Contains(output.String(), "example removed.") {
		t.Fatalf("failure presentation is not truthful: %q", output.String())
	}
}

func TestRunStorageFailureDoesNotReportToolRemoved(t *testing.T) {
	storage := filepath.Join(t.TempDir(), "example")
	if err := os.MkdirAll(storage, 0755); err != nil {
		t.Fatal(err)
	}
	tool := fakeTool{name: "example", storage: storage}
	installed := map[string]bool{"example": true}
	deps := testWorkflowDependencies(t, []string{"example"}, map[string]tools.Tool{"example": tool}, installed)
	deps.removeAll = func(string) error { return errors.New("permission denied") }
	var output bytes.Buffer
	err := run(context.Background(), Streams{In: strings.NewReader(""), Out: &output}, Options{Tool: "example", AssumeYes: true}, deps)
	if err == nil || !strings.Contains(err.Error(), "clean example storage") {
		t.Fatalf("error = %v, want storage context", err)
	}
	if !strings.Contains(output.String(), "cleaning example storage... ✗") || strings.Contains(output.String(), "example removed.") {
		t.Fatalf("storage failure presentation is not truthful: %q", output.String())
	}
}

func TestRunCancellationDuringProviderRemovalIsExpected(t *testing.T) {
	tool := fakeTool{name: "example"}
	installed := map[string]bool{"example": true}
	deps := testWorkflowDependencies(t, []string{"example"}, map[string]tools.Tool{"example": tool}, installed)
	deps.uninstall = func(context.Context, tools.Tool, *environment.Environment) error { return context.Canceled }
	var output bytes.Buffer
	if err := run(context.Background(), Streams{In: strings.NewReader(""), Out: &output}, Options{Tool: "example", AssumeYes: true}, deps); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "uninstalling example... cancelled\nuninstall cancelled.") {
		t.Fatalf("cancellation output = %q", output.String())
	}
	if strings.Contains(output.String(), "✓") || strings.Contains(output.String(), "✗") || strings.Contains(output.String(), "removed.") {
		t.Fatalf("cancellation implies completion: %q", output.String())
	}
}

func TestRunDeclinedConfirmationDoesNotMutate(t *testing.T) {
	tool := fakeTool{name: "example"}
	installed := map[string]bool{"example": true}
	deps := testWorkflowDependencies(t, []string{"example"}, map[string]tools.Tool{"example": tool}, installed)
	called := false
	deps.uninstall = func(context.Context, tools.Tool, *environment.Environment) error { called = true; return nil }
	var output bytes.Buffer
	if err := run(context.Background(), Streams{In: strings.NewReader("n\n"), Out: &output}, Options{Tool: "example"}, deps); err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("declined uninstall mutated the tool")
	}
	if !strings.Contains(output.String(), "continue uninstall? (y/N): uninstall cancelled.") {
		t.Fatalf("decline output = %q", output.String())
	}
}
