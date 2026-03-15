package uninstall

import (
	"os"
	"path/filepath"

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

		// uninstall tool
		err := tool.Uninstall(env)
		if err != nil {
			return err
		}

		// cleanup mdev storage directory
		storagePath := filepath.Join(env.DataRoot, name)

		if _, err := os.Stat(storagePath); err == nil {
			printer.Info("Cleaning storage: " + storagePath)

			err := os.RemoveAll(storagePath)
			if err != nil {
				return err
			}
		}

		printer.Success(name + " removed")
	}

	return nil
}
