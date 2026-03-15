package doctor

import (
	"time"

	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Fix attempts to resolve issues detected by doctor.
// Currently focuses on system prerequisites.
func Fix() {

	var missing []prerequisites.Prerequisite

	for _, p := range prerequisites.List() {
		if !p.Check() {
			missing = append(missing, p)
		}
	}

	if len(missing) == 0 {
		printer.Success("Nothing to fix")
		return
	}

	printer.Section("Fixing system prerequisites")

	for _, m := range missing {
		printer.Info(m.Name())
	}

	if !interactive.AskYesNo("Install missing prerequisites?") {
		printer.Info("Aborted")
		return
	}

	printer.Section("Fixing system")

	startAll := time.Now()

	for _, m := range missing {

		printer.Blank()
		printer.Info("Installing " + m.Name())

		start := time.Now()

		err := m.Install()
		if err != nil {
			printer.Fail(m.Name() + " installation failed")
			continue
		}

		elapsed := time.Since(start).Round(time.Second)

		printer.Success(m.Name() + " installed")
		printer.Indent(1, "time: "+elapsed.String())
	}

	total := time.Since(startAll).Round(time.Second)

	printer.Blank()
	printer.Section("Summary")
	printer.Info("Total fix time: " + total.String())
}
