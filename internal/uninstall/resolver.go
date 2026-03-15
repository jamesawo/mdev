package uninstall

import "github.com/jamesawo/mdev/internal/tools"

// FindDependents returns tools that depend on the given tool.
func FindDependents(name string) []string {

	var dependents []string

	for _, t := range tools.List() {

		for _, dep := range t.Dependencies() {
			if dep == name {
				dependents = append(dependents, t.Name())
			}
		}
	}

	return dependents
}
