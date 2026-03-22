package doctor

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Report represents the result of a doctor run.
type Report struct {
	System      []Check
	Environment []Check
	Tools       []ToolCheck
}

// Check represents a generic system/environment check.
type Check struct {
	Name   string
	Status bool
	Detail string
}

// ToolCheck represents the status of a development tool.
type ToolCheck struct {
	Name         string
	Installed    bool
	Dependencies []string
}

// Execute runs the doctor command flow.
func Execute(isFixFlag bool) {

	if isFixFlag {
		Fix()
		return
	}

	report, err := Run()
	if err != nil {
		printer.Fail(messages.DoctorFailed)
		return
	}

	renderReport(report)
}

// Run executes all doctor checks and returns a structured report.
func Run() (*Report, error) {

	report := &Report{}

	// Phase 1: system prerequisites
	sys := checkSystemPrerequisites()
	report.System = sys

	// Phase 2: environment
	envChecks := checkEnvironment()
	report.Environment = envChecks

	// Phase 3: tools
	toolChecks := checkTools()
	report.Tools = toolChecks

	return report, nil
}

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
