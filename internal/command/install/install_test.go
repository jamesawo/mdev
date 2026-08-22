package install

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
)

type fakeTool struct {
	name  string
	deps  []string
	calls *[]string
	fail  string
}

func (f *fakeTool) Name() string                               { return f.name }
func (f *fakeTool) Description() string                        { return f.name }
func (f *fakeTool) Dependencies() []string                     { return f.deps }
func (f *fakeTool) IsInstalled(*environment.Environment) bool  { return false }
func (f *fakeTool) Install(*environment.Environment) error     { return f.record("install") }
func (f *fakeTool) Configure(*environment.Environment) error   { return f.record("configure") }
func (f *fakeTool) Verify(*environment.Environment) error      { return f.record("verify") }
func (f *fakeTool) Uninstall(*environment.Environment) error   { return nil }
func (f *fakeTool) StorageDir(*environment.Environment) string { return "" }
func (f *fakeTool) record(phase string) error {
	*f.calls = append(*f.calls, phase+":"+f.name)
	if f.fail == phase {
		return errors.New("failed")
	}
	return nil
}

func testDependencies(all []tools.Tool, plan []tools.Tool, installed map[string]bool) workflowDependencies {
	byName := make(map[string]tools.Tool, len(all))
	for _, tool := range all {
		byName[tool.Name()] = tool
	}
	return workflowDependencies{
		loadEnvironment: func() (*environment.Environment, bool, error) { return environment.New("/unused"), true, nil },
		validateStorage: func(*environment.Environment) error { return nil },
		listTools:       func() []tools.Tool { return all },
		getTool:         func(name string) (tools.Tool, bool) { tool, ok := byName[name]; return tool, ok },
		resolve:         func([]tools.Tool) ([]tools.Tool, error) { return plan, nil },
		status:          func(tool tools.Tool, _ *environment.Environment) (bool, error) { return installed[tool.Name()], nil },
		selectTools:     func([]choice) ([]string, error) { return nil, nil },
		newReporter:     func(out io.Writer) progressReporter { return newTextProgressReporter(out) },
	}
}

func TestRunExecutesDependencyFirstWithStableProgress(t *testing.T) {
	var calls []string
	dependency := &fakeTool{name: "dependency", calls: &calls}
	requested := &fakeTool{name: "requested", deps: []string{"dependency"}, calls: &calls}
	deps := testDependencies([]tools.Tool{requested, dependency}, []tools.Tool{dependency, requested}, map[string]bool{})
	var output bytes.Buffer
	err := run(context.Background(), Streams{In: strings.NewReader(""), Out: &output}, Options{Tool: "requested", AssumeYes: true}, deps)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"install:dependency", "configure:dependency", "verify:dependency", "install:requested", "configure:requested", "verify:requested"}
	if strings.Join(calls, ",") != strings.Join(want, ",") {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
	for _, text := range []string{"Install plan", "dependency", "requested", "installed"} {
		if !strings.Contains(output.String(), text) {
			t.Fatalf("output omits %q:\n%s", text, output.String())
		}
	}
	if strings.Contains(output.String(), "\x1b[") {
		t.Fatalf("output contains terminal control sequences: %q", output.String())
	}
}

func TestRunStopsAfterPhaseFailure(t *testing.T) {
	for _, phase := range []string{"install", "configure", "verify"} {
		t.Run(phase, func(t *testing.T) {
			var calls []string
			tool := &fakeTool{name: "example", calls: &calls, fail: phase}
			deps := testDependencies([]tools.Tool{tool}, []tools.Tool{tool}, map[string]bool{})
			err := run(context.Background(), Streams{In: strings.NewReader(""), Out: &bytes.Buffer{}}, Options{Tool: tool.name, AssumeYes: true}, deps)
			if err == nil || !strings.Contains(err.Error(), phase) {
				t.Fatalf("error = %v, want phase context", err)
			}
			if calls[len(calls)-1] != phase+":"+tool.name {
				t.Fatalf("calls = %v", calls)
			}
		})
	}
}

func TestRunReportsAlreadyInstalledWithoutVersion(t *testing.T) {
	var calls []string
	tool := &fakeTool{name: "example", calls: &calls}
	deps := testDependencies([]tools.Tool{tool}, []tools.Tool{tool}, map[string]bool{"example": true})
	var output bytes.Buffer
	if err := run(context.Background(), Streams{In: strings.NewReader(""), Out: &output}, Options{Tool: "example", AssumeYes: true}, deps); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 0 {
		t.Fatalf("calls = %v", calls)
	}
	if !strings.Contains(output.String(), "already installed") || !strings.Contains(output.String(), "uninstall") {
		t.Fatalf("output = %q", output.String())
	}
	if strings.Contains(output.String(), "0.0") {
		t.Fatalf("output includes a version: %q", output.String())
	}
}

func TestRunRejectsUnknownToolBeforeMutation(t *testing.T) {
	deps := testDependencies(nil, nil, nil)
	err := run(context.Background(), Streams{In: strings.NewReader(""), Out: &bytes.Buffer{}}, Options{Tool: "missing", AssumeYes: true}, deps)
	if err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("error = %v", err)
	}
}

func TestRunCancellationDoesNotMutate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	deps := testDependencies(nil, nil, nil)
	var output bytes.Buffer
	if err := run(ctx, Streams{In: strings.NewReader(""), Out: &output}, Options{}, deps); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "cancelled") {
		t.Fatalf("output = %q", output.String())
	}
}
