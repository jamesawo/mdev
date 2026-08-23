package uninstall

import (
	"errors"
	"os"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

var askConfirmation = confirmation.Ask

func Run(env *environment.Environment, name string) error {

	// ---- Resolve uninstall plan (dependency-aware) ----
	resolved, err := BuildPlan(name)
	if err != nil {
		return err
	}
	plan, err := installedPlan(resolved, func(name string) (bool, error) {
		tool, ok := tools.Get(name)
		if !ok {
			return false, nil
		}
		return tools.InstallationStatus(tool, env)
	})
	if err != nil {
		return err
	}
	if len(plan) == 0 {
		printer.Info(messages.UninstallNotInstalled(name))
		return nil
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

		if !askConfirmation(messages.UninstallRemoveDependentsQuestion) {
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
	var directories []string
	for _, name := range plan {
		tool, ok := tools.Get(name)
		if !ok {
			continue
		}
		path := tool.StorageDir(env)
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			directories = append(directories, path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if len(directories) > 0 {
		printer.Section(messages.UninstallDirectoriesToRemove)
		for _, path := range directories {
			printer.Info(path)
		}
	}

	if !askConfirmation(messages.UninstallContinueQuestion) {
		printer.Info(messages.UninstallCancelled)
		return nil
	}

	// ---- Execute uninstall ----
	return execute(env, plan)
}

func installedPlan(resolved []string, status func(string) (bool, error)) ([]string, error) {
	plan := make([]string, 0, len(resolved))
	for _, name := range resolved {
		installed, err := status(name)
		if err != nil {
			return nil, err
		}
		if installed {
			plan = append(plan, name)
		}
	}
	return plan, nil
}
