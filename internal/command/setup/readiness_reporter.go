package setup

import (
	"fmt"
	"io"

	"github.com/jamesawo/mdev/internal/readiness"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type readinessReporter struct{ out io.Writer }

func newReadinessReporter(out io.Writer) readiness.Reporter { return &readinessReporter{out: out} }
func (r *readinessReporter) Checking(name string) error {
	_, err := fmt.Fprintf(r.out, messages.SetupReadinessChecking, name)
	return err
}
func (r *readinessReporter) Checked(item readiness.Item) error {
	if item.Ready() {
		_, err := fmt.Fprintf(r.out, messages.SetupReadinessReady, item.Prerequisite.Name())
		return err
	}
	_, err := fmt.Fprintf(r.out, messages.SetupReadinessNeeds, item.Prerequisite.Name(), item.RemediationDescription())
	return err
}
func (r *readinessReporter) Remediating(name, description string) error {
	_, err := fmt.Fprintf(r.out, messages.SetupReadinessRemediating, description, name)
	return err
}
func (r *readinessReporter) Verified(item readiness.Item) error {
	_, err := fmt.Fprintf(r.out, messages.SetupReadinessVerified, item.Prerequisite.Name())
	return err
}
