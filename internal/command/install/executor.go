package install

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

func execute(env *environment.Environment, plan []tools.Tool) error {

	// ---- Install Plan Preview ----
	printer.Section(messages.ToolsInstallPlan)

	for _, t := range plan {
		printer.Info(t.Name())
	}

	ok := interactive.AskYesNo(messages.ToolsContinueInstallQuestion)
	if !ok {
		printer.Info(messages.ToolsInstallCancelled)
		printer.Blank()
		return nil
	}

	// ---- Installation Execution ----
	printer.Section(messages.ToolsInstallingStart)

	for _, tool := range plan {

		if tool.IsInstalled(env) {
			printer.Success(tool.Name() + " " + messages.ToolsAlreadyInstalled)
			continue
		}

		printer.Info(messages.Installing + " " + tool.Name())

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

		printer.Success(tool.Name() + " " + messages.Installed)
	}

	return nil
}
