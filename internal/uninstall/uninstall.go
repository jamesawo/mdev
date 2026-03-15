package uninstall

import (
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Run orchestrates the uninstall workflow.
func Run(env *environment.Environment, name string) error {

	// ---- Resolve uninstall plan (dependency-aware) ----
	plan, err := BuildPlan(name)
	if err != nil {
		return err
	}

	// If more than one tool appears in the plan,
	// it means other tools depend on the target.
	if len(plan) > 1 {

		printer.Section("Dependency warning")

		printer.Info(name + " is required by:")

		// dependents are everything except the last item (target)
		for _, dep := range plan[:len(plan)-1] {
			printer.Info("  " + dep)
		}

		if !interactive.AskYesNo("Remove dependent tools first?") {
			printer.Info("Cancelled.")
			return nil
		}
	}

	// ---- Show uninstall plan ----
	printer.Section("Uninstall plan")

	for _, tool := range plan {
		printer.Info(tool)
	}

	// Show directories that will be removed
	printer.Section("Directories to be removed")

	for _, tool := range plan {
		path := StoragePath(env, tool)
		printer.Info(path)
	}

	if !interactive.AskYesNo("Continue uninstall?") {
		printer.Info("Cancelled.")
		return nil
	}

	if !interactive.AskYesNo("Continue uninstall?") {
		printer.Info("Cancelled.")
		return nil
	}

	// ---- Execute uninstall ----
	return execute(env, plan)
}
