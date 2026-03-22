package doctor

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
)

// checkTools reports installation status of tools.
func checkTools(reporter Reporter) []ToolCheck {

	results := []ToolCheck{}

	env, _ := environment.FromConfig()

	for _, t := range tools.List() {

		result := ToolCheck{
			Name:         t.Name(),
			Installed:    t.IsInstalled(env),
			Dependencies: t.Dependencies(),
		}
		results = append(results, result)
		if reporter != nil {
			reporter.ToolCheck(result)
		}
	}

	return results
}
