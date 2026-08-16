package uninstall

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

func Run(env *environment.Environment, name string) error {

	// ---- Resolve uninstall plan (dependency-aware) ----
	plan, err := BuildPlan(name)
	if err != nil {
		return err
	}

	// If more than one tool appears in the plan,
	// it means other tools depend on the target.
	if len(plan) > 1 {

		printer.Section(messages.UninstallDependencyWarning)

		printer.Info(messages.UninstallRequiredBy(name))

		// dependents are everything except the last item (target)
		for _, dep := range plan[:len(plan)-1] {
			printer.Info("  " + dep)
		}

		if !confirmation.Ask(messages.UninstallRemoveDependentsQuestion) {
			printer.Info(messages.UninstallCancelled)
			return nil
		}
	}

	// ---- Show uninstall plan ----
	printer.Section(messages.UninstallPlan)

	for _, tool := range plan {
		printer.Info(tool)
	}

	// Show directories that will be removed
	printer.Section(messages.UninstallDirectoriesToRemove)

	for _, tool := range plan {
		path := StoragePath(env, tool)
		printer.Info(path)
	}

	if !confirmation.Ask(messages.UninstallContinueQuestion) {
		printer.Info(messages.UninstallCancelled)
		return nil
	}

	if !confirmation.Ask(messages.UninstallContinueQuestion) {
		printer.Info(messages.UninstallCancelled)
		return nil
	}

	// ---- Execute uninstall ----
	return execute(env, plan)
}
