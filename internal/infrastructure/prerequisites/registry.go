package prerequisites

import (
	"fmt"
	"sort"

	"github.com/jamesawo/mdev/internal/ui/messages"
)

var registry []Prerequisite

// Register adds a prerequisite to the global registry.
func Register(p Prerequisite) {
	registry = append(registry, p)
}

// List returns all registered prerequisites.
func List() []Prerequisite {
	result := append([]Prerequisite(nil), registry...)
	sort.Slice(result, func(i, j int) bool { return result[i].Name() < result[j].Name() })
	return result
}

// Get returns a prerequisite by exact canonical name.
func Get(name string) (Prerequisite, bool) {
	for _, prerequisite := range registry {
		if prerequisite.Name() == name {
			return prerequisite, true
		}
	}
	return nil, false
}

// Resolve returns a deterministic dependency-first prerequisite closure.
func Resolve(selected []Prerequisite) ([]Prerequisite, error) {
	roots := append([]Prerequisite(nil), selected...)
	sort.Slice(roots, func(i, j int) bool { return roots[i].Name() < roots[j].Name() })
	visited := map[string]bool{}
	active := map[string]bool{}
	var result []Prerequisite
	var visit func(Prerequisite) error
	visit = func(prerequisite Prerequisite) error {
		name := prerequisite.Name()
		if active[name] {
			return fmt.Errorf(messages.ReadinessDependencyCycleError, name)
		}
		if visited[name] {
			return nil
		}
		active[name] = true
		var names []string
		if provider, ok := prerequisite.(DependencyProvider); ok {
			names = append(names, provider.PrerequisiteDependencies()...)
			sort.Strings(names)
		}
		for _, dependencyName := range names {
			dependency, ok := Get(dependencyName)
			if !ok {
				return fmt.Errorf(messages.ReadinessUnknownDependencyError, name, dependencyName)
			}
			if err := visit(dependency); err != nil {
				return err
			}
		}
		active[name] = false
		visited[name] = true
		result = append(result, prerequisite)
		return nil
	}
	for _, prerequisite := range roots {
		if err := visit(prerequisite); err != nil {
			return nil, err
		}
	}
	return result, nil
}
