package doctor

import (
	"context"

	"github.com/jamesawo/mdev/internal/ui/messages"
)

// Run executes all doctor checks and returns a structured report.
func Run(reporter Reporter) (*Report, error) {
	return RunContext(context.Background(), reporter)
}

// RunContext executes doctor checks with cancellation.
func RunContext(ctx context.Context, reporter Reporter) (*Report, error) {
	if reporter == nil {
		reporter = noopReporter{}
	}

	report := &Report{}

	// Phase 1: system prerequisites
	reporter.StartSection(messages.DoctorSystem)
	sys, err := checkSystemPrerequisites(ctx, reporter)
	if err != nil {
		return nil, err
	}
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
