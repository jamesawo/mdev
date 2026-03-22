package doctor

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
