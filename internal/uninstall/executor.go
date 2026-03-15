package uninstall

import (
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

func execute(env *environment.Environment, names []string) error {

	printer.Section("Uninstalling tools")

	for _, name := range names {

		tool, ok := tools.Get(name)
		if !ok {
			continue
		}

		if !tool.IsInstalled(env) {
			printer.Info(name + " not installed")
			continue
		}

		printer.Info("Removing " + name)

		err := tool.Uninstall(env)
		if err != nil {
			return err
		}

		printer.Success(name + " removed")
	}

	return nil
}
