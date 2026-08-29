package uninstall

import (
	"fmt"
	"io"

	"github.com/jamesawo/mdev/internal/ui/messages"
)

type progressReporter interface {
	DependencyWarning(string, []string) error
	Plan([]string) error
	Directories([]string) error
	Started(string, string) error
	Succeeded(string, string) error
	Failed(string, string) error
	Completed([]string) error
	NotInstalled(string) error
	SystemRequirement(string) error
	Cancelled() error
}

type textProgressReporter struct {
	out       io.Writer
	phaseOpen bool
}

func newTextProgressReporter(out io.Writer) *textProgressReporter {
	return &textProgressReporter{out: out}
}

func (r *textProgressReporter) DependencyWarning(name string, dependants []string) error {
	if _, err := fmt.Fprintln(r.out, messages.UninstallDependencyWarning); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(r.out, messages.UninstallRequiredBy, name); err != nil {
		return err
	}
	for _, dependant := range dependants {
		if _, err := fmt.Fprintf(r.out, messages.UninstallDependantItem, dependant); err != nil {
			return err
		}
	}
	return nil
}

func (r *textProgressReporter) Plan(plan []string) error {
	if _, err := fmt.Fprintln(r.out, messages.UninstallPlan); err != nil {
		return err
	}
	for _, name := range plan {
		if _, err := fmt.Fprintf(r.out, messages.UninstallItem, name); err != nil {
			return err
		}
	}
	return nil
}

func (r *textProgressReporter) Directories(directories []string) error {
	if len(directories) == 0 {
		return nil
	}
	if _, err := fmt.Fprintln(r.out, messages.UninstallDirectoriesToRemove); err != nil {
		return err
	}
	for _, path := range directories {
		if _, err := fmt.Fprintf(r.out, messages.UninstallItem, path); err != nil {
			return err
		}
	}
	return nil
}

func (r *textProgressReporter) Started(name, phase string) error {
	var err error
	if phase == messages.UninstallPhaseStorage {
		_, err = fmt.Fprintf(r.out, messages.UninstallStorageProgress, name)
	} else {
		_, err = fmt.Fprintf(r.out, messages.UninstallToolProgress, name)
	}
	if err == nil {
		r.phaseOpen = true
	}
	return err
}

func (r *textProgressReporter) Succeeded(_, _ string) error {
	r.phaseOpen = false
	_, err := fmt.Fprint(r.out, messages.UninstallPhaseSucceeded)
	return err
}

func (r *textProgressReporter) Failed(_, _ string) error {
	r.phaseOpen = false
	_, err := fmt.Fprint(r.out, messages.UninstallPhaseFailed)
	return err
}

func (r *textProgressReporter) Completed(names []string) error {
	for _, name := range names {
		if _, err := fmt.Fprintf(r.out, messages.UninstallRemoved, name); err != nil {
			return err
		}
	}
	return nil
}

func (r *textProgressReporter) NotInstalled(name string) error {
	_, err := fmt.Fprintf(r.out, messages.UninstallNotInstalled, name)
	return err
}

func (r *textProgressReporter) SystemRequirement(name string) error {
	_, err := fmt.Fprintf(r.out, messages.UninstallSystemRequirement, name)
	return err
}

func (r *textProgressReporter) Cancelled() error {
	if r.phaseOpen {
		r.phaseOpen = false
		if _, err := fmt.Fprint(r.out, messages.UninstallPhaseCancelled); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(r.out, messages.UninstallCancelled)
	return err
}
