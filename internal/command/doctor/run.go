package doctor

import "github.com/jamesawo/mdev/internal/ui/messages"

// Run executes all doctor checks and returns a structured report.
func Run(reporter Reporter) (*Report, error) {
	if reporter == nil {
		reporter = noopReporter{}
	}

	report := &Report{}

	// Phase 1: system prerequisites
	reporter.StartSection(messages.DoctorSystem)
	sys := checkSystemPrerequisites(reporter)
	report.System = sys

	// Phase 2: environment
	reporter.StartSection(messages.DoctorEnvironment)
	envChecks := checkEnvironment(reporter)
	report.Environment = envChecks

	// Phase 3: tools
	reporter.StartSection(messages.DoctorTools)
	toolChecks := checkTools(reporter)
	report.Tools = toolChecks

	return report, nil
}
