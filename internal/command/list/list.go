package list

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Run prints all tools supported by mdev and their installation status.
//
// The command loads the current environment configuration (if available)
// and checks each registered tool to determine whether it is installed.
// Tools that are installed are displayed with a success marker,
// while tools not installed are displayed with a failure marker.
func Run() {

	env, _ := environment.FromConfig()

	printer.Section("Available tools")

	for _, t := range tools.List() {

		name := t.Name()

		// If environment exists and the tool is installed,
		// mark it as installed in the output.
		if env != nil && t.IsInstalled(env) {
			printer.Success(name + " (installed)")
			continue
		}

		// Otherwise show it as available but not installed.
		printer.Fail(name)
	}
}
