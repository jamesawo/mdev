package doctor

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// renderReport renders the doctor report.
func renderReport(report *Report) {

	printer.Section(messages.System)

	for _, s := range report.System {
		if s.Status {
			printer.Success(s.Name)
		} else {
			printer.Fail(fmt.Sprintf("%s %s", messages.Missing, s.Name))
		}
	}

	printer.Section(messages.Environment)

	for _, e := range report.Environment {
		if e.Status {
			if e.Detail != "" {
				printer.Success(e.Name + ": " + e.Detail)
			} else {
				printer.Success(e.Name)
			}
		} else {
			printer.Fail(messages.DoctorNotConfigured(e.Name))
			printer.Indent(2, messages.Run+" mdev install")
		}
	}

	printer.Section(messages.DoctorTools)

	for _, t := range report.Tools {
		if t.Installed {
			printer.Success(t.Name)
			continue
		}

		printer.Fail(t.Name)
	}

	printer.Section(messages.DoctorNextSteps)

	printer.Indent(1, messages.DoctorInstallIndividual)
	for _, t := range report.Tools {
		if !t.Installed {
			printer.Indent(2, messages.DoctorInstallTool(t.Name))
		}
	}

	printer.Blank()
	printer.Info(messages.DoctorInstallEverything)
	printer.Indent(2, "mdev install --all")

	printer.Blank()
	printer.Info(messages.DoctorFixHint)
	printer.Indent(2, "mdev doctor --fix")
	printer.Blank()
}
