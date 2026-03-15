package doctor

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
