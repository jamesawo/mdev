package tools

import (
	"fmt"
	"sort"
)

// ErrUnknownDependency creates the canonical missing-dependency error.
func ErrUnknownDependency(name string) error {
	return fmt.Errorf("unknown dependency %s", name)
}

// ResolveOrder returns a deterministic dependency-first plan for all tools.
func ResolveOrder() ([]Tool, error) {
	return ResolveDependencies(List())
}

// ResolveDependencies returns a deterministic dependency-first closure for the
// selected tools without validating unrelated registry entries.
func ResolveDependencies(selected []Tool) ([]Tool, error) {
	roots := append([]Tool(nil), selected...)
	sort.Slice(roots, func(i, j int) bool { return roots[i].Name() < roots[j].Name() })
	visited := map[string]bool{}
	stack := map[string]bool{}
	order := make([]Tool, 0)
	var visit func(Tool) error
	visit = func(tool Tool) error {
		name := tool.Name()
		if stack[name] {
			return fmt.Errorf("dependency cycle detected at %s", name)
		}
		if visited[name] {
			return nil
		}
		stack[name] = true
		dependencies := append([]string(nil), tool.Dependencies()...)
		sort.Strings(dependencies)
		for _, dependency := range dependencies {
			dependencyTool, ok := Get(dependency)
			if !ok {
				return fmt.Errorf("%s requires %s: %w", name, dependency, ErrUnknownDependency(dependency))
			}
			if err := visit(dependencyTool); err != nil {
				return err
			}
		}
		stack[name] = false
		visited[name] = true
		order = append(order, tool)
		return nil
	}
	for _, tool := range roots {
		if err := visit(tool); err != nil {
			return nil, err
		}
	}
	return order, nil
}
