package doctor

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
)

// checkTools reports installation status of tools.
func checkTools(reporter Reporter) []ToolCheck {

	results := []ToolCheck{}

	env, _ := environment.FromConfig()

	registered := tools.List()
	nameWidth := 0
	for _, tool := range registered {
		if len(tool.Name()) > nameWidth {
			nameWidth = len(tool.Name())
		}
	}

	for _, t := range registered {
		reporter.StartCheck(t.Name())

		result := ToolCheck{
			Name:         t.Name(),
			Installed:    t.IsInstalled(env),
			Dependencies: t.Dependencies(),
			NameWidth:    nameWidth,
		}
		results = append(results, result)
		reporter.ToolCheck(result)
	}

	return results
}
