package doctor

import "github.com/jamesawo/mdev/internal/system"

// checkSystem validates system prerequisites.
func checkSystem() []Check {

	results := []Check{}

	for _, p := range system.List() {

		ok := p.Check()

		results = append(results, Check{
			Name:   p.Name(),
			Status: ok,
		})
	}

	return results
}
