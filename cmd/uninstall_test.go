package cmd

import (
	"context"
	"errors"
	"testing"

	commanduninstall "github.com/jamesawo/mdev/internal/command/uninstall"
)

func TestUninstallCommandDelegatesOptionsAndStreams(t *testing.T) {
	var got commanduninstall.Options
	runUninstall = func(_ context.Context, streams commanduninstall.Streams, options commanduninstall.Options) error {
		got = options
		if streams.In == nil || streams.Out == nil || streams.Err == nil {
			t.Fatal("missing command streams")
		}
		return nil
	}
	t.Cleanup(func() { runUninstall = defaultRunUninstall })
	confirmAll = true
	t.Cleanup(func() { confirmAll = false })
	if err := uninstallCmd.RunE(uninstallCmd, []string{"podman"}); err != nil {
		t.Fatal(err)
	}
	if got.Tool != "podman" || !got.AssumeYes {
		t.Fatalf("options = %#v", got)
	}
}

func TestUninstallCommandReturnsWorkflowError(t *testing.T) {
	want := errors.New("uninstall failed")
	runUninstall = func(context.Context, commanduninstall.Streams, commanduninstall.Options) error { return want }
	t.Cleanup(func() { runUninstall = defaultRunUninstall })
	if err := uninstallCmd.RunE(uninstallCmd, []string{"podman"}); !errors.Is(err, want) {
		t.Fatalf("error = %v", err)
	}
}

func TestUninstallCommandRequiresExactlyOneTool(t *testing.T) {
	if err := uninstallCmd.Args(uninstallCmd, nil); err == nil {
		t.Fatal("accepted missing tool")
	}
	if err := uninstallCmd.Args(uninstallCmd, []string{"one", "two"}); err == nil {
		t.Fatal("accepted extra tool")
	}
}
