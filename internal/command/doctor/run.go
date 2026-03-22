package doctor

import "github.com/jamesawo/mdev/internal/ui/messages"

// Run executes all doctor checks and returns a structured report.
func Run(reporter Reporter) (*Report, error) {

	report := &Report{}

	// Phase 1: system prerequisites
	if reporter != nil {
		reporter.StartSection(messages.System)
	}
	sys := checkSystemPrerequisites(reporter)
	report.System = sys

	// Phase 2: environment
	if reporter != nil {
		reporter.StartSection(messages.Environment)
	}
	envChecks := checkEnvironment(reporter)
	report.Environment = envChecks

	// Phase 3: tools
	if reporter != nil {
		reporter.StartSection(messages.DoctorTools)
	}
	toolChecks := checkTools(reporter)
	report.Tools = toolChecks

	return report, nil
}
