package tools

import (
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

type dependencyTool struct {
	name         string
	dependencies []string
}

func (t dependencyTool) Name() string                             { return t.name }
func (t dependencyTool) Description() string                      { return t.name }
func (t dependencyTool) Dependencies() []string                   { return t.dependencies }
func (dependencyTool) IsInstalled(*environment.Environment) bool  { return false }
func (dependencyTool) Install(*environment.Environment) error     { return nil }
func (dependencyTool) Configure(*environment.Environment) error   { return nil }
func (dependencyTool) Verify(*environment.Environment) error      { return nil }
func (dependencyTool) Uninstall(*environment.Environment) error   { return nil }
func (dependencyTool) StorageDir(*environment.Environment) string { return "" }

func TestResolveDependenciesIsScopedDeterministicAndDependencyFirst(t *testing.T) {
	previous := registry
	registry = map[string]Tool{}
	t.Cleanup(func() { registry = previous })
	Register(dependencyTool{name: "root", dependencies: []string{"z", "a"}})
	Register(dependencyTool{name: "a"})
	Register(dependencyTool{name: "z"})
	Register(dependencyTool{name: "unrelated", dependencies: []string{"missing"}})
	root, _ := Get("root")
	plan, err := ResolveDependencies([]Tool{root})
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(plan))
	for _, tool := range plan {
		names = append(names, tool.Name())
	}
	if strings.Join(names, ",") != "a,z,root" {
		t.Fatalf("plan = %v", names)
	}
}

func TestResolveDependenciesReportsCycle(t *testing.T) {
	previous := registry
	registry = map[string]Tool{}
	t.Cleanup(func() { registry = previous })
	Register(dependencyTool{name: "a", dependencies: []string{"b"}})
	Register(dependencyTool{name: "b", dependencies: []string{"a"}})
	a, _ := Get("a")
	_, err := ResolveDependencies([]Tool{a})
	if err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("error = %v", err)
	}
}

func TestSystemPrerequisitesAreDeduplicatedAndSorted(t *testing.T) {
	plan := []Tool{
		systemPrerequisiteTool{dependencyTool{name: "one"}, []string{"git", "brew"}},
		systemPrerequisiteTool{dependencyTool{name: "two"}, []string{"brew", "bash"}},
	}
	got := SystemPrerequisites(plan)
	if strings.Join(got, ",") != "bash,brew,git" {
		t.Fatalf("requirements = %v", got)
	}
}

type systemPrerequisiteTool struct {
	dependencyTool
	requirements []string
}

func (t systemPrerequisiteTool) SystemPrerequisites() []string { return t.requirements }
