package doctor

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

type noopReporter struct{}

func (noopReporter) StartSection(string)    {}
func (noopReporter) StartCheck(string)      {}
func (noopReporter) SystemCheck(Check)      {}
func (noopReporter) EnvironmentCheck(Check) {}
func (noopReporter) ToolCheck(ToolCheck)    {}

type progressReporter struct{}

func (r *progressReporter) StartSection(title string) {
	printer.Section(title)
}

func (r *progressReporter) StartCheck(name string) {
	printer.Info(fmt.Sprintf(messages.DoctorCheckingProgress, name))
}

func (r *progressReporter) SystemCheck(result Check) {
	if result.Status {
		printer.Success(result.Name)
		return
	}

	printer.Fail(fmt.Sprintf(messages.DoctorMissingCheck, messages.DoctorMissing, result.Name))
}

func (r *progressReporter) EnvironmentCheck(result Check) {
	if result.Status {
		if result.Detail != "" {
			printer.Success(fmt.Sprintf(messages.DoctorCheckDetail, result.Name, result.Detail))
		} else {
			printer.Success(result.Name)
		}
		return
	}

	printer.Fail(messages.DoctorEnvironmentNotSetUp)
	printer.Indent(2, messages.CommonRun+" "+messages.SetupCommand)
}

func (r *progressReporter) ToolCheck(result ToolCheck) {
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
	printer.Indent(2, messages.DoctorInstallAllCommand)

	printer.Blank()
	printer.Info(messages.DoctorFixHint)
	printer.Indent(2, messages.DoctorFixCommand)
	printer.Blank()
}
