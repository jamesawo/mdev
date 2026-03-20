package doctor

import (
	"time"

	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
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
		printer.Success(messages.DoctorNothingToFix)
		printer.Blank()
		return
	}

	printer.Section(messages.DoctorFixingPrerequisites)

	for _, m := range missing {
		printer.Info(m.Name())
	}

	if !interactive.AskYesNo(messages.DoctorInstallMissingPrerequisites) {
		printer.Info(messages.Aborted)
		return
	}

	printer.Section(messages.DoctorFixingSystem)

	startAll := time.Now()

	for _, m := range missing {

		printer.Blank()
		printer.Info(messages.Installing + " " + m.Name())

		start := time.Now()

		err := m.Install()
		if err != nil {
			printer.Fail(messages.DoctorInstallationFailed(m.Name()))
			continue
		}

		elapsed := time.Since(start).Round(time.Second)

		printer.Success(m.Name() + " " + messages.Installed)
		printer.Indent(1, messages.DoctorTimeElapsed(elapsed.String()))
		printer.Blank()
	}

	total := time.Since(startAll).Round(time.Second)

	printer.Blank()
	printer.Section(messages.DoctorFixSummary)
	printer.Info(messages.DoctorTotalFixTime(total.String()))
	printer.Blank()
}
