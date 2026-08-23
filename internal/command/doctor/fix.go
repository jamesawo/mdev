package doctor

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jamesawo/mdev/internal/readiness"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

// Fix attempts to resolve issues detected by doctor.
func Fix() error { return FixContext(context.Background()) }

// FixContext checks and repairs system readiness using process streams.
func FixContext(ctx context.Context) error {
	return FixContextWithStreams(ctx, os.Stdin, os.Stdout)
}

// FixContextWithStreams checks, confirms, remediates, and verifies system readiness.
func FixContextWithStreams(ctx context.Context, in io.Reader, out io.Writer) error {
	reporter := &doctorFixReporter{out: out}
	items, err := readiness.CheckAll(ctx, reporter)
	if err != nil {
		return err
	}
	unready := readiness.Unready(items)
	if len(unready) == 0 {
		_, err := fmt.Fprintln(out, messages.DoctorNothingToFix)
		return err
	}
	if err := reporter.Summary(unready); err != nil {
		return err
	}
	if !confirmation.New(in, out, false).AskDefaultNo(messages.DoctorInstallMissingPrerequisites) {
		_, err := fmt.Fprintln(out, messages.CommonAborted)
		return err
	}
	if err := readiness.Remediate(ctx, items, reporter); err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, messages.DoctorSystemReady)
	return err
}

type doctorFixReporter struct {
	out      io.Writer
	checking bool
	fixing   bool
}

func (r *doctorFixReporter) Checking(string) error {
	if r.checking {
		return nil
	}
	r.checking = true
	_, err := fmt.Fprintln(r.out, messages.DoctorCheckingSystem)
	return err
}

func (r *doctorFixReporter) Checked(item readiness.Item) error {
	if !item.Ready() {
		return nil
	}
	_, err := fmt.Fprintf(r.out, messages.DoctorReadyResult, item.Prerequisite.Name())
	return err
}

func (r *doctorFixReporter) Summary(items []readiness.Item) error {
	if _, err := fmt.Fprintf(r.out, "\n%s\n", messages.DoctorNeedsAttention); err != nil {
		return err
	}
	for _, item := range items {
		detail := item.Detail
		if detail == "" {
			detail = item.RemediationDescription()
		}
		if _, err := fmt.Fprintf(r.out, messages.DoctorAttentionItem, item.Prerequisite.Name(), strings.ToLower(detail)); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(r.out, "\n%s\n", messages.DoctorSystemChanges(len(items)))
	return err
}

func (r *doctorFixReporter) Remediating(name, _ string) error {
	if !r.fixing {
		r.fixing = true
		if _, err := fmt.Fprintf(r.out, "\n%s\n", messages.DoctorFixingSystem); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(r.out, messages.DoctorFixingItem, name)
	return err
}

func (r *doctorFixReporter) Verified(item readiness.Item) error {
	_, err := fmt.Fprintf(r.out, messages.DoctorReadyResult, item.Prerequisite.Name())
	return err
}
