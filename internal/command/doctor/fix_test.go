package doctor

import (
	"bytes"
	"context"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
	"github.com/jamesawo/mdev/internal/readiness"
)

func TestDoctorFixReporterShowsResultsWithoutNarratingEveryCheck(t *testing.T) {
	var output bytes.Buffer
	reporter := &doctorFixReporter{out: &output}
	items := []readiness.Item{
		{Prerequisite: doctorReporterPrerequisite{name: "curl"}, State: prerequisites.StateReady},
		{Prerequisite: doctorReporterPrerequisite{name: "brew"}, State: prerequisites.StateReady},
	}
	for _, item := range items {
		if err := reporter.Checking(item.Prerequisite.Name()); err != nil {
			t.Fatal(err)
		}
		if err := reporter.Checked(item); err != nil {
			t.Fatal(err)
		}
	}
	want := "checking system requirements...\n✓ curl\n✓ brew\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestDoctorFixReporterSeparatesAttentionFromCompletedChecks(t *testing.T) {
	var output bytes.Buffer
	reporter := &doctorFixReporter{out: &output}
	items := []readiness.Item{
		{Prerequisite: doctorReporterPrerequisite{name: "brew"}, State: prerequisites.StateMissing, Detail: "Homebrew is required"},
		{Prerequisite: doctorReporterPrerequisite{name: "bash"}, State: prerequisites.StateOutdated, Detail: "Bash 4 or newer is required"},
	}
	if err := reporter.Summary(items); err != nil {
		t.Fatal(err)
	}
	want := "\nneeds attention\n" +
		"  brew\n    homebrew is required\n" +
		"  bash\n    bash 4 or newer is required\n" +
		"\n2 system changes can be fixed by mdev.\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestDoctorFixReporterShowsRealRemediationProgress(t *testing.T) {
	var output bytes.Buffer
	reporter := &doctorFixReporter{out: &output}
	item := readiness.Item{Prerequisite: doctorReporterPrerequisite{name: "brew"}, State: prerequisites.StateReady}
	if err := reporter.Remediating("brew", "install Homebrew"); err != nil {
		t.Fatal(err)
	}
	if err := reporter.Verified(item); err != nil {
		t.Fatal(err)
	}
	want := "\nfixing system requirements...\nfixing brew...\n✓ brew\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

type doctorReporterPrerequisite struct{ name string }

func (p doctorReporterPrerequisite) Name() string { return p.name }
func (doctorReporterPrerequisite) Check() bool    { return false }
func (doctorReporterPrerequisite) Install() error { return nil }
func (doctorReporterPrerequisite) RemediationDescription() string {
	return "install requirement"
}
func (doctorReporterPrerequisite) RemediateContext(context.Context) error { return nil }
