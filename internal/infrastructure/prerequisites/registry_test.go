package prerequisites

import (
	"reflect"
	"strings"
	"testing"
)

type registryPrerequisite struct {
	name         string
	dependencies []string
}

func (p registryPrerequisite) Name() string                       { return p.name }
func (registryPrerequisite) Check() bool                          { return true }
func (registryPrerequisite) Install() error                       { return nil }
func (p registryPrerequisite) PrerequisiteDependencies() []string { return p.dependencies }

func TestResolveIsDeterministicDependencyFirst(t *testing.T) {
	previous := registry
	registry = nil
	t.Cleanup(func() { registry = previous })
	Register(registryPrerequisite{name: "root", dependencies: []string{"z", "a"}})
	Register(registryPrerequisite{name: "z"})
	Register(registryPrerequisite{name: "a"})
	root, _ := Get("root")
	plan, err := Resolve([]Prerequisite{root})
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, prerequisite := range plan {
		names = append(names, prerequisite.Name())
	}
	if !reflect.DeepEqual(names, []string{"a", "z", "root"}) {
		t.Fatalf("plan = %v", names)
	}
}

func TestResolveReportsUnknownDependency(t *testing.T) {
	previous := registry
	registry = nil
	t.Cleanup(func() { registry = previous })
	root := registryPrerequisite{name: "root", dependencies: []string{"missing"}}
	Register(root)
	_, err := Resolve([]Prerequisite{root})
	if err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("error = %v", err)
	}
}
