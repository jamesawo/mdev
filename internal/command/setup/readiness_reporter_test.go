package setup

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
	"github.com/jamesawo/mdev/internal/readiness"
)

func TestReadinessReporterPresentsOneProgressLineAndReadableSummary(t *testing.T) {
	var output bytes.Buffer
	reporter := newReadinessReporter(&output).(*readinessReporter)
	items := []readiness.Item{
		{Prerequisite: reporterPrerequisite{name: "curl", ready: true}, State: prerequisites.StateReady},
		{Prerequisite: reporterPrerequisite{name: "xcode-cli", remediation: "install Xcode Command Line Tools"}, State: prerequisites.StateMissing},
		{Prerequisite: reporterPrerequisite{name: "brew", remediation: "install Homebrew"}, State: prerequisites.StateMissing},
		{Prerequisite: reporterPrerequisite{name: "git", ready: true}, State: prerequisites.StateReady},
	}
	for _, item := range items {
		if err := reporter.Checking(item.Prerequisite.Name()); err != nil {
			t.Fatal(err)
		}
		if err := reporter.Checked(item); err != nil {
			t.Fatal(err)
		}
	}
	if err := reporter.Summary(items); err != nil {
		t.Fatal(err)
	}
	want := "checking system requirements...\n" +
		"✓ curl\n" +
		"✓ git\n" +
		"\nsystem changes required\n\n" +
		"  xcode-cli\n" +
		"    install Xcode Command Line Tools\n\n" +
		"  brew\n" +
		"    install Homebrew\n\n" +
		"2 system changes are required to finish setup.\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
	if strings.Count(output.String(), "checking") != 1 {
		t.Fatalf("checking output was repeated: %q", output.String())
	}
}

func TestReadinessReporterExplainsDeferredInstallation(t *testing.T) {
	var output bytes.Buffer
	reporter := newReadinessReporter(&output).(*readinessReporter)
	if err := reporter.Deferred("installation started.", "finish it, then rerun setup."); err != nil {
		t.Fatal(err)
	}
	want := "installation started.\nfinish it, then rerun setup.\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestReadinessReporterShowsRealRemediationProgress(t *testing.T) {
	var output bytes.Buffer
	reporter := newReadinessReporter(&output).(*readinessReporter)
	brew := readiness.Item{Prerequisite: reporterPrerequisite{name: "brew"}, State: prerequisites.StateReady}
	bash := readiness.Item{Prerequisite: reporterPrerequisite{name: "bash"}, State: prerequisites.StateReady}
	if err := reporter.Remediating("brew", "install Homebrew"); err != nil {
		t.Fatal(err)
	}
	if err := reporter.Verified(brew); err != nil {
		t.Fatal(err)
	}
	if err := reporter.Remediating("bash", "install Bash"); err != nil {
		t.Fatal(err)
	}
	if err := reporter.Verified(bash); err != nil {
		t.Fatal(err)
	}
	want := "\ninstalling system requirements...\n" +
		"installing brew...\n" +
		"✓ brew\n" +
		"installing bash...\n" +
		"✓ bash\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

type reporterPrerequisite struct {
	name        string
	remediation string
	ready       bool
}

func (p reporterPrerequisite) Name() string { return p.name }
func (p reporterPrerequisite) Check() bool  { return p.ready }
func (reporterPrerequisite) Install() error { return nil }
func (p reporterPrerequisite) RemediationDescription() string {
	return p.remediation
}
func (reporterPrerequisite) RemediateContext(context.Context) error { return nil }
