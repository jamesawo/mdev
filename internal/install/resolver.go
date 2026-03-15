package install

import "github.com/jamesawo/mdev/internal/tools"

// resolveSelection expands dependencies
// and returns tools in correct install order.
func resolveSelection(selected []tools.Tool) ([]tools.Tool, error) {

	graph := map[string]bool{}

	for _, t := range selected {
		graph[t.Name()] = true

		for _, dep := range t.Dependencies() {
			graph[dep] = true
		}
	}

	var names []string

	for n := range graph {
		names = append(names, n)
	}

	return tools.ResolveSubset(names)
}
