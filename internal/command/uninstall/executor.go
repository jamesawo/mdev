package uninstall

import (
	"errors"
	"os"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

func execute(env *environment.Environment, names []string) error {

	printer.Section(messages.UninstallTools)

	for _, name := range names {

		tool, ok := tools.Get(name)
		if !ok {
			continue
		}

		installed, err := tools.InstallationStatus(tool, env)
		if err != nil {
			return err
		}
		if !installed {
			printer.Info(messages.UninstallNotInstalled(name))
			continue
		}

		printer.Info(messages.UninstallRemoving(name))

		// uninstall tool
		err = tool.Uninstall(env)
		if err != nil {
			return err
		}

		// cleanup mdev storage directory
		storagePath := tool.StorageDir(env)

		if storagePath != "" {
			if _, err := os.Stat(storagePath); err == nil {
				printer.Info(messages.UninstallCleaningStorage(storagePath))
				if err := os.RemoveAll(storagePath); err != nil {
					return err
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}

		printer.Success(messages.UninstallRemoved(name))
	}

	return nil
}
