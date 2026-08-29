package install

import (
	"bytes"
	"strings"
	"testing"
)

func TestTextReporterMarksPhaseOnlyAfterSuccess(t *testing.T) {
	var output bytes.Buffer
	reporter := newTextProgressReporter(&output)
	if err := reporter.Started("java", "installing"); err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != "installing java... " {
		t.Fatalf("started output = %q", got)
	}
	if err := reporter.Succeeded("java", "installing"); err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != "installing java... ✓\n" {
		t.Fatalf("succeeded output = %q", got)
	}
}

func TestTextReporterFailureDoesNotRenderSuccess(t *testing.T) {
	var output bytes.Buffer
	reporter := newTextProgressReporter(&output)
	if err := reporter.Started("java", "configuring"); err != nil {
		t.Fatal(err)
	}
	if err := reporter.Failed("java", "configuring"); err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != "configuring java... ✗\n" {
		t.Fatalf("failure output = %q", got)
	}
	if strings.Contains(output.String(), "installed.") {
		t.Fatalf("failure output implies final success: %q", output.String())
	}
}

func TestTextReporterCancellationClosesActivePhase(t *testing.T) {
	var output bytes.Buffer
	reporter := newTextProgressReporter(&output)
	if err := reporter.Started("java", "installing"); err != nil {
		t.Fatal(err)
	}
	if err := reporter.Cancelled(); err != nil {
		t.Fatal(err)
	}
	want := "installing java... cancelled\ninstallation cancelled.\n"
	if got := output.String(); got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
