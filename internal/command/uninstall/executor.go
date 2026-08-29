package uninstall

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/subprocess"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

func execute(ctx context.Context, env *environment.Environment, names []string, reporter progressReporter, deps workflowDependencies) error {
	completed := make([]string, 0, len(names))
	for _, name := range names {
		if ctx.Err() != nil {
			return reporter.Cancelled()
		}
		tool, ok := deps.getTool(name)
		if !ok {
			continue
		}
		installed, err := deps.status(name, env)
		if err != nil {
			return err
		}
		if !installed {
			continue
		}

		if err := reporter.Started(name, messages.UninstallPhaseTool); err != nil {
			return err
		}
		phaseContext := subprocess.WithManagedOutput(ctx)
		if err := deps.uninstall(phaseContext, tool, env); err != nil {
			if ctx.Err() != nil || errors.Is(err, context.Canceled) {
				return reporter.Cancelled()
			}
			if reportErr := reporter.Failed(name, messages.UninstallPhaseTool); reportErr != nil {
				return errors.Join(fmt.Errorf(messages.UninstallToolError, name, err), reportErr)
			}
			return fmt.Errorf(messages.UninstallToolError, name, err)
		}
		if err := reporter.Succeeded(name, messages.UninstallPhaseTool); err != nil {
			return err
		}

		storagePath := tool.StorageDir(env)
		if storagePath != "" {
			if _, err := deps.stat(storagePath); err == nil {
				if err := reporter.Started(name, messages.UninstallPhaseStorage); err != nil {
					return err
				}
				if err := deps.removeAll(storagePath); err != nil {
					if reportErr := reporter.Failed(name, messages.UninstallPhaseStorage); reportErr != nil {
						return errors.Join(fmt.Errorf(messages.UninstallStorageError, name, err), reportErr)
					}
					return fmt.Errorf(messages.UninstallStorageError, name, err)
				}
				if err := reporter.Succeeded(name, messages.UninstallPhaseStorage); err != nil {
					return err
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
		completed = append(completed, name)
	}
	return reporter.Completed(completed)
}
