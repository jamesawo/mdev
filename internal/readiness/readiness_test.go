package readiness

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
)

type fakePrerequisite struct {
	name         string
	state        prerequisites.State
	calls        *[]string
	remediateErr error
	checkErr     error
}

func (p *fakePrerequisite) Name() string   { return p.name }
func (p *fakePrerequisite) Check() bool    { return p.state == prerequisites.StateReady }
func (p *fakePrerequisite) Install() error { return nil }
func (p *fakePrerequisite) Readiness(context.Context) (prerequisites.State, string, error) {
	return p.state, "", p.checkErr
}
func (p *fakePrerequisite) RemediationDescription() string { return "prepare " + p.name }
func (p *fakePrerequisite) RemediateContext(context.Context) error {
	*p.calls = append(*p.calls, "remediate:"+p.name)
	if p.remediateErr != nil {
		return p.remediateErr
	}
	p.state = prerequisites.StateReady
	return nil
}
func (p *fakePrerequisite) VerifyContext(context.Context) error {
	*p.calls = append(*p.calls, "verify:"+p.name)
	return nil
}

func TestRemediatePerformsAndVerifiesApprovedWork(t *testing.T) {
	var calls []string
	p := &fakePrerequisite{name: "example", state: prerequisites.StateMissing, calls: &calls}
	items, err := Check(context.Background(), []prerequisites.Prerequisite{p}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := Remediate(context.Background(), items, nil); err != nil {
		t.Fatal(err)
	}
	want := []string{"remediate:example", "verify:example"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestRemediatePreservesFailureAndStops(t *testing.T) {
	var calls []string
	wantErr := errors.New("failed")
	p := &fakePrerequisite{name: "example", state: prerequisites.StateMissing, calls: &calls, remediateErr: wantErr}
	err := Remediate(context.Background(), []Item{{Prerequisite: p, State: prerequisites.StateMissing}}, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v", err)
	}
	if len(calls) != 1 {
		t.Fatalf("calls = %v", calls)
	}
}

type manualPrerequisite struct{ name string }

func (p manualPrerequisite) Name() string { return p.name }
func (manualPrerequisite) Check() bool    { return false }
func (manualPrerequisite) Install() error { return nil }

func TestRemediateRefusesManualOnlyPrerequisite(t *testing.T) {
	err := Remediate(context.Background(), []Item{{Prerequisite: manualPrerequisite{"manual"}, State: prerequisites.StateMissing}}, nil)
	if err == nil {
		t.Fatal("manual prerequisite was treated as remediable")
	}
}

func TestCheckPreservesUnknownStatusFailure(t *testing.T) {
	wantErr := errors.New("status unavailable")
	calls := []string{}
	p := &fakePrerequisite{name: "example", state: prerequisites.StateBroken, calls: &calls, checkErr: wantErr}
	_, err := Check(context.Background(), []prerequisites.Prerequisite{p}, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v", err)
	}
}

func TestRemediateHonorsCancellationBeforeMutation(t *testing.T) {
	var calls []string
	p := &fakePrerequisite{name: "example", state: prerequisites.StateMissing, calls: &calls}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Remediate(ctx, []Item{{Prerequisite: p, State: prerequisites.StateMissing}}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v", err)
	}
	if len(calls) != 0 {
		t.Fatalf("calls = %v", calls)
	}
}

type deferredPrerequisite struct {
	fakePrerequisite
}

func (*deferredPrerequisite) RemediationStartedMessage() string      { return "started" }
func (*deferredPrerequisite) RemediationContinuationMessage() string { return "continue externally" }

type recordingDeferredReporter struct {
	noopReporter
	started      string
	continuation string
}

func (r *recordingDeferredReporter) Deferred(started, continuation string) error {
	r.started = started
	r.continuation = continuation
	return nil
}

func TestRemediateDefersVerificationUntilExternalInstallationCompletes(t *testing.T) {
	var calls []string
	p := &deferredPrerequisite{fakePrerequisite: fakePrerequisite{
		name: "xcode-cli", state: prerequisites.StateMissing, calls: &calls,
	}}
	reporter := &recordingDeferredReporter{}
	err := Remediate(context.Background(), []Item{{Prerequisite: p, State: prerequisites.StateMissing}}, reporter)
	var pending *PendingRemediationError
	if !errors.As(err, &pending) || pending.Prerequisite != "xcode-cli" {
		t.Fatalf("error = %v", err)
	}
	if !reflect.DeepEqual(calls, []string{"remediate:xcode-cli"}) {
		t.Fatalf("calls = %v", calls)
	}
	if reporter.started != "started" || reporter.continuation != "continue externally" {
		t.Fatalf("deferred output = %q, %q", reporter.started, reporter.continuation)
	}
}
