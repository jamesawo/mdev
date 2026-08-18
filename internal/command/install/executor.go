package install

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

func execute(env *environment.Environment, plan []tools.Tool) error {

	// ---- Install Plan Preview ----
	printer.Section(messages.InstallPlan)

	for _, t := range plan {
		printer.Info(t.Name())
	}

	ok := confirmation.Ask(messages.InstallContinueQuestion)
	if !ok {
		printer.Info(messages.InstallCancelled)
		printer.Blank()
		return nil
	}

	// ---- Installation Execution ----
	printer.Section(messages.InstallStart)

	for _, tool := range plan {

		if tool.IsInstalled(env) {
			printer.Success(tool.Name() + " " + messages.InstallAlreadyInstalled)
			continue
		}

		printer.Info(messages.CommonInstalling + " " + tool.Name())

		err := tool.Install(env)
		if err != nil {
			return err
		}

		err = tool.Configure(env)
		if err != nil {
			return err
		}

		err = tool.Verify(env)
		if err != nil {
			return err
		}

		printer.Success(tool.Name() + " " + messages.CommonInstalled)
	}

	return nil
}
