package subprocess

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestManagedRunSuppressesSuccessfulOutput(t *testing.T) {
	ctx := WithManagedOutput(context.Background())
	if err := Run(ctx, "sh", "-c", "printf 'provider stdout'; printf 'provider stderr' >&2"); err != nil {
		t.Fatal(err)
	}
}

func TestManagedRunRetainsFailureDiagnosticsAndCause(t *testing.T) {
	ctx := WithManagedOutput(context.Background())
	err := Run(ctx, "sh", "-c", "printf 'download chatter\\n'; printf 'useful failure' >&2; exit 7")
	if err == nil || !strings.Contains(err.Error(), "useful failure") {
		t.Fatalf("error = %v, want captured diagnostic", err)
	}
	if strings.Contains(err.Error(), "download chatter") {
		t.Fatalf("error includes successful stdout despite useful stderr: %v", err)
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 7 {
		t.Fatalf("error = %v, want preserved exit status 7", err)
	}
}

func TestManagedRunUsesStdoutWhenFailureHasNoStderr(t *testing.T) {
	ctx := WithManagedOutput(context.Background())
	err := Run(ctx, "sh", "-c", "printf 'stdout-only failure'; exit 1")
	if err == nil || !strings.Contains(err.Error(), "stdout-only failure") {
		t.Fatalf("error = %v, want stdout diagnostic", err)
	}
}

func TestManagedRunBoundsFailureDiagnostics(t *testing.T) {
	ctx := WithManagedOutput(context.Background())
	err := Run(ctx, "sh", "-c", "i=0; while [ $i -lt 100 ]; do printf 'provider line %s xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\n' $i; i=$((i+1)); done; exit 1")
	if err == nil {
		t.Fatal("command unexpectedly succeeded")
	}
	if len(err.Error()) > maxDiagnosticBytes+100 {
		t.Fatalf("diagnostic is not bounded: %d bytes", len(err.Error()))
	}
	if strings.Contains(err.Error(), "provider line 0 ") || !strings.Contains(err.Error(), "provider line 99 ") {
		t.Fatalf("diagnostic does not retain the useful tail: %s", err)
	}
}

func TestCheckSuppressesSuccessAndReportsFailure(t *testing.T) {
	if err := Check(context.Background(), "sh", "-c", "printf 'version output'"); err != nil {
		t.Fatal(err)
	}
	err := Check(context.Background(), "sh", "-c", "printf 'verification failed' >&2; exit 1")
	if err == nil || !strings.Contains(err.Error(), "verification failed") {
		t.Fatalf("error = %v, want verification diagnostic", err)
	}
}

func TestManagedRunPreservesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(WithManagedOutput(context.Background()))
	cancel()
	err := Run(ctx, "sh", "-c", "sleep 10")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context cancellation", err)
	}
}
