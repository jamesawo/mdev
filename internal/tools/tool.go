package tools

import (
	"sort"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

type Tool interface {
	Name() string
	Description() string

	Dependencies() []string

	IsInstalled(env *environment.Environment) bool
	Install(env *environment.Environment) error
	Configure(env *environment.Environment) error
	Verify(env *environment.Environment) error
	Uninstall(env *environment.Environment) error

	StorageDir(env *environment.Environment) string
}

// SystemPrerequisiteProvider optionally declares host capabilities required by a tool.
type SystemPrerequisiteProvider interface {
	SystemPrerequisites() []string
}

// SystemPrerequisites returns the deduplicated host requirements for a tool plan.
func SystemPrerequisites(plan []Tool) []string {
	seen := map[string]bool{}
	var result []string
	for _, tool := range plan {
		provider, ok := tool.(SystemPrerequisiteProvider)
		if !ok {
			continue
		}
		for _, name := range provider.SystemPrerequisites() {
			if !seen[name] {
				seen[name] = true
				result = append(result, name)
			}
		}
	}
	sort.Strings(result)
	return result
}
