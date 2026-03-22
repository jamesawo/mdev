package doctor

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

type progressReporter struct{}

func (progressReporter) StartSection(title string) {
	printer.Section(title)
}

func (progressReporter) SystemCheck(result Check) {
	if result.Status {
		printer.Success(result.Name)
		return
	}

	printer.Fail(fmt.Sprintf("%s %s", messages.Missing, result.Name))
}

func (progressReporter) EnvironmentCheck(result Check) {
	if result.Status {
		if result.Detail != "" {
			printer.Success(result.Name + ": " + result.Detail)
		} else {
			printer.Success(result.Name)
		}
		return
	}

	printer.Fail(messages.DoctorNotConfigured(result.Name))
	printer.Indent(2, messages.Run+" mdev install")
}

func (progressReporter) ToolCheck(result ToolCheck) {
	if result.Installed {
		printer.Success(result.Name)
		return
	}

	printer.Fail(result.Name)
}

// renderSummary renders the final doctor summary after streaming progress.
func renderSummary(report *Report) {

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
