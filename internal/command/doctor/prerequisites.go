package doctor

import (
	"context"

	"github.com/jamesawo/mdev/internal/readiness"
)

// checkSystemPrerequisites checks all registered prerequisites,
// prints their status, and installs missing ones if the user agrees.
func checkSystemPrerequisites(ctx context.Context, reporter Reporter) ([]Check, error) {
	adapter := &doctorReadinessReporter{reporter: reporter}
	_, err := readiness.CheckAll(ctx, adapter)
	if err != nil {
		return nil, err
	}
	return adapter.checks, nil
}

type doctorReadinessReporter struct {
	reporter Reporter
	checks   []Check
}

func (r *doctorReadinessReporter) Checking(name string) error {
	r.reporter.StartCheck(name)
	return nil
}
func (r *doctorReadinessReporter) Checked(item readiness.Item) error {
	result := Check{Name: item.Prerequisite.Name(), Status: item.Ready(), Detail: item.Detail}
	r.checks = append(r.checks, result)
	r.reporter.SystemCheck(result)
	return nil
}
func (*doctorReadinessReporter) Remediating(string, string) error { return nil }
func (*doctorReadinessReporter) Verified(readiness.Item) error    { return nil }
