package install

import (
	"fmt"
	"io"

	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

// progressReporter receives install lifecycle events and separates orchestration
// from how those events are presented to the user.
type progressReporter interface {
	Plan([]tools.Tool) error
	Started(string, string) error
	Succeeded(string, string) error
	Failed(string, string) error
	Completed(string) error
	AlreadyInstalled(string, bool) error
	Cancelled() error
	NoSelection() error
}

// textProgressReporter renders lifecycle events as stable plain text.
type textProgressReporter struct {
	out       io.Writer
	phaseOpen bool
}

// newTextProgressReporter creates the default writer-backed progress reporter.
func newTextProgressReporter(out io.Writer) *textProgressReporter {
	return &textProgressReporter{out: out}
}

// Plan writes the dependency-first plan before confirmation.
func (r *textProgressReporter) Plan(plan []tools.Tool) error {
	if _, err := fmt.Fprintln(r.out, messages.InstallPlan); err != nil {
		return err
	}
	for _, tool := range plan {
		if _, err := fmt.Fprintf(r.out, messages.InstallPlanItem, tool.Name()); err != nil {
			return err
		}
	}
	return nil
}

// Started writes a lifecycle boundary immediately before work begins.
func (r *textProgressReporter) Started(name, phase string) error {
	_, err := fmt.Fprintf(r.out, messages.InstallProgress, phase, name)
	if err == nil {
		r.phaseOpen = true
	}
	return err
}

// Succeeded closes the active lifecycle line only after its work succeeds.
func (r *textProgressReporter) Succeeded(_, _ string) error {
	r.phaseOpen = false
	_, err := fmt.Fprint(r.out, messages.InstallPhaseSucceeded)
	return err
}

// Failed closes the active lifecycle line without implying phase success.
func (r *textProgressReporter) Failed(_, _ string) error {
	r.phaseOpen = false
	_, err := fmt.Fprint(r.out, messages.InstallPhaseFailed)
	return err
}

// Completed reports a tool only after verification succeeds.
func (r *textProgressReporter) Completed(name string) error {
	_, err := fmt.Fprintf(r.out, messages.InstallCompleted, name)
	return err
}

// AlreadyInstalled reports a skipped tool and optionally its recovery command.
func (r *textProgressReporter) AlreadyInstalled(name string, recovery bool) error {
	if _, err := fmt.Fprintf(r.out, messages.InstallAlreadyInstalled, name); err != nil {
		return err
	}
	if recovery {
		_, err := fmt.Fprintf(r.out, messages.InstallUninstallHint, name)
		return err
	}
	return nil
}

// Cancelled writes the concise successful-cancellation message.
func (r *textProgressReporter) Cancelled() error {
	if r.phaseOpen {
		r.phaseOpen = false
		if _, err := fmt.Fprint(r.out, messages.InstallPhaseCancelled); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(r.out, messages.InstallCancelled)
	return err
}

// NoSelection reports an interactive request that selected no tools.
func (r *textProgressReporter) NoSelection() error {
	_, err := fmt.Fprintln(r.out, messages.InstallNoSelection)
	return err
}
