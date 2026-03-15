package uninstall

import (
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

func Run(env *environment.Environment, name string) error {

	dependents := FindDependents(name)

	var plan []string

	// If other tools depend on this one
	if len(dependents) > 0 {

		printer.Section("Dependency warning")

		printer.Info(name + " is required by:")

		for _, d := range dependents {
			printer.Info("  " + d)
		}

		if !interactive.AskYesNo("Remove dependent tools first?") {
			printer.Info("Cancelled.")
			return nil
		}

		plan = append(plan, dependents...)
	}

	plan = append(plan, name)

	// ---- uninstall plan preview ----
	printer.Section("Uninstall plan")

	for _, p := range plan {
		printer.Info(p)
	}

	if !interactive.AskYesNo("Continue uninstall?") {
		printer.Info("Cancelled.")
		return nil
	}

	return execute(env, plan)
}
