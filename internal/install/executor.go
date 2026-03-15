package install

import (
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

func execute(env *environment.Environment, plan []tools.Tool) error {

	// ---- Install Plan Preview ----
	printer.Section("Install plan")

	for _, t := range plan {
		printer.Info(t.Name())
	}

	ok := interactive.AskYesNo("Continue installation?")
	if !ok {
		printer.Info("Installation cancelled.")
		return nil
	}

	// ---- Installation Execution ----
	printer.Section("Installing tools")

	for _, tool := range plan {

		if tool.IsInstalled(env) {
			printer.Success(tool.Name() + " already installed")
			continue
		}

		printer.Info("Installing " + tool.Name())

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

		printer.Success(tool.Name() + " installed")
	}

	return nil
}
