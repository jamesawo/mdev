package doctor

import (
	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
)

// checkSystemPrerequisites checks all registered prerequisites,
// prints their status, and installs missing ones if the user agrees.
func checkSystemPrerequisites() []Check {

	var checks []Check

	for _, p := range prerequisites.List() {

		ok := p.Check()

		checks = append(checks, Check{
			Name:   p.Name(),
			Status: ok,
			Detail: "",
		})
	}

	return checks
}
