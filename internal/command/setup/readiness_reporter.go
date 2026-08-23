package setup

import (
	"fmt"
	"io"

	"github.com/jamesawo/mdev/internal/readiness"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type readinessReporter struct {
	out        io.Writer
	checking   bool
	installing bool
}

func newReadinessReporter(out io.Writer) readiness.Reporter { return &readinessReporter{out: out} }
func (r *readinessReporter) Checking(string) error {
	if r.checking {
		return nil
	}
	r.checking = true
	_, err := fmt.Fprint(r.out, messages.SetupReadinessChecking)
	return err
}
func (r *readinessReporter) Checked(item readiness.Item) error {
	if item.Ready() {
		_, err := fmt.Fprintf(r.out, messages.SetupReadinessReady, item.Prerequisite.Name())
		return err
	}
	return nil
}
func (r *readinessReporter) Summary(items []readiness.Item) error {
	unready := readiness.Unready(items)
	if len(unready) == 0 {
		return nil
	}
	if _, err := fmt.Fprint(r.out, messages.SetupReadinessRequired); err != nil {
		return err
	}
	for _, item := range unready {
		if _, err := fmt.Fprintf(r.out, messages.SetupReadinessChange, item.Prerequisite.Name(), item.RemediationDescription()); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(r.out, messages.SetupReadinessChanges(len(unready)))
	return err
}
func (r *readinessReporter) Remediating(name, _ string) error {
	if !r.installing {
		r.installing = true
		if _, err := fmt.Fprint(r.out, messages.SetupReadinessInstalling); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(r.out, messages.SetupReadinessRemediating, name)
	return err
}
func (r *readinessReporter) Verified(item readiness.Item) error {
	_, err := fmt.Fprintf(r.out, messages.SetupReadinessVerified, item.Prerequisite.Name())
	return err
}
func (r *readinessReporter) Deferred(started, continuation string) error {
	_, err := fmt.Fprintf(r.out, messages.SetupReadinessDeferred, started, continuation)
	return err
}
