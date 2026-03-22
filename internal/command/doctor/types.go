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

// Reporter streams progress updates while doctor checks are running.
type Reporter interface {
	StartSection(title string)
	SystemCheck(result Check)
	EnvironmentCheck(result Check)
	ToolCheck(result ToolCheck)
}
