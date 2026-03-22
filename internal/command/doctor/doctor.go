package doctor

import (
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Execute runs the doctor command flow.
func Execute(isFixFlag bool) {

	if isFixFlag {
		Fix()
		return
	}

	report, err := Run(progressReporter{})

	if err != nil {
		printer.Fail(messages.DoctorFailed)
		return
	}

	renderSummary(report)
}
