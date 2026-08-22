package doctor

import (
	"context"
	"fmt"

	"github.com/jamesawo/mdev/internal/readiness"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Fix attempts to resolve issues detected by doctor.
// Currently focuses on system prerequisites.
func Fix() error { return FixContext(context.Background()) }

// FixContext checks, confirms, remediates, and verifies system readiness.
func FixContext(ctx context.Context) error {
	reporter := doctorFixReporter{}
	items, err := readiness.CheckAll(ctx, reporter)
	if err != nil {
		return err
	}
	unready := readiness.Unready(items)
	if len(unready) == 0 {
		printer.Success(messages.DoctorNothingToFix)
		printer.Blank()
		return nil
	}

	printer.Section(messages.DoctorFixingPrerequisites)

	for _, item := range unready {
		printer.Info(item.Prerequisite.Name())
	}

	if !confirmation.Ask(messages.DoctorInstallMissingPrerequisites) {
		printer.Info(messages.CommonAborted)
		return nil
	}

	printer.Section(messages.DoctorFixingSystem)

	if err := readiness.Remediate(ctx, items, reporter); err != nil {
		return err
	}
	printer.Success(messages.DoctorNothingToFix)
	return nil
}

type doctorFixReporter struct{}

func (doctorFixReporter) Checking(name string) error {
	printer.Info(fmt.Sprintf(messages.DoctorChecking, name))
	return nil
}
func (doctorFixReporter) Checked(item readiness.Item) error {
	if item.Ready() {
		printer.Success(item.Prerequisite.Name())
	} else {
		printer.Fail(fmt.Sprintf(messages.DoctorReadinessState, item.Prerequisite.Name(), item.State))
	}
	return nil
}
func (doctorFixReporter) Remediating(name, description string) error {
	printer.Info(fmt.Sprintf(messages.DoctorRemediating, description, name))
	return nil
}
func (doctorFixReporter) Verified(item readiness.Item) error {
	printer.Success(fmt.Sprintf(messages.DoctorVerified, item.Prerequisite.Name()))
	return nil
}
