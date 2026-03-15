package doctor

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
)

// checkTools reports installation status of tools.
func checkTools() []ToolCheck {

	results := []ToolCheck{}

	env, _ := environment.FromConfig()

	for _, t := range tools.List() {

		results = append(results, ToolCheck{
			Name:         t.Name(),
			Installed:    t.IsInstalled(env),
			Dependencies: t.Dependencies(),
		})
	}

	return results
}
