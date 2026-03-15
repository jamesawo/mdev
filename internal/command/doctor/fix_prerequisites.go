package doctor

import (
	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// FixMissingPrerequisites installs prerequisites that failed checks.
func FixMissingPrerequisites(checks []Check) {

	var missing []prerequisites.Prerequisite

	for _, p := range prerequisites.List() {

		if !p.Check() {
			missing = append(missing, p)
		}
	}

	if len(missing) == 0 {
		return
	}

	printer.Section("Missing prerequisites")

	for _, m := range missing {
		printer.Info(m.Name())
	}

	if !interactive.AskYesNo("Install missing prerequisites?") {
		return
	}

	printer.Section("Installing prerequisites")

	for _, m := range missing {

		printer.Info("Installing " + m.Name())

		if err := m.Install(); err != nil {
			printer.Fail(m.Name() + " installation failed")
			printer.Blank()
			continue
		}

		printer.Success(m.Name() + " installed")
	}
}
