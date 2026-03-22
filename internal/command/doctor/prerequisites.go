package doctor

import (
	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
)

// checkSystemPrerequisites checks all registered prerequisites,
// prints their status, and installs missing ones if the user agrees.
func checkSystemPrerequisites(reporter Reporter) []Check {

	var checks []Check

	for _, p := range prerequisites.List() {
		reporter.StartCheck(p.Name())

		ok := p.Check()

		result := Check{
			Name:   p.Name(),
			Status: ok,
			Detail: "",
		}
		checks = append(checks, result)
		reporter.SystemCheck(result)
	}

	return checks
}
