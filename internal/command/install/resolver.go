package install

import "github.com/jamesawo/mdev/internal/tools"

// resolveSelection expands dependencies recursively
// and returns tools in the correct installation order.
func resolveSelection(selected []tools.Tool) ([]tools.Tool, error) {

	visited := map[string]bool{}

	// recursively collect dependencies
	var collect func(t tools.Tool)

	collect = func(t tools.Tool) {

		if visited[t.Name()] {
			return
		}

		visited[t.Name()] = true

		for _, dep := range t.Dependencies() {

			d, ok := tools.Get(dep)
			if !ok {
				continue
			}

			collect(d)
		}
	}

	for _, t := range selected {
		collect(t)
	}

	var names []string

	for n := range visited {
		names = append(names, n)
	}

	return tools.ResolveSubset(names)
}
