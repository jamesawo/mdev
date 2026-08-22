package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	commandinstall "github.com/jamesawo/mdev/internal/command/install"
)

func TestInstallCommandDelegatesOptionsAndStreams(t *testing.T) {
	var got commandinstall.Options
	runInstall = func(_ context.Context, streams commandinstall.Streams, options commandinstall.Options) error {
		got = options
		if streams.In == nil || streams.Out == nil || streams.Err == nil {
			t.Fatal("missing command streams")
		}
		return nil
	}
	t.Cleanup(func() { runInstall = defaultRunInstall })
	installAll = false
	confirmAll = true
	t.Cleanup(func() { installAll = false; confirmAll = false })
	if err := installCmd.RunE(installCmd, []string{"gradle"}); err != nil {
		t.Fatal(err)
	}
	if got.Tool != "gradle" || got.All || !got.AssumeYes {
		t.Fatalf("options = %#v", got)
	}
}

func TestInstallCommandReturnsWorkflowError(t *testing.T) {
	want := errors.New("install failed")
	runInstall = func(context.Context, commandinstall.Streams, commandinstall.Options) error { return want }
	t.Cleanup(func() { runInstall = defaultRunInstall })
	if err := installCmd.RunE(installCmd, nil); !errors.Is(err, want) {
		t.Fatalf("error = %v", err)
	}
}

func TestInstallCommandRejectsAllWithTool(t *testing.T) {
	installAll = true
	t.Cleanup(func() { installAll = false })
	if err := installCmd.PreRunE(installCmd, []string{"gradle"}); err == nil {
		t.Fatal("accepted --all with a tool")
	}
}

func TestInstallCommandRejectsExtraArguments(t *testing.T) {
	if err := installCmd.Args(installCmd, []string{"one", "two"}); err == nil {
		t.Fatal("accepted extra arguments")
	}
}

func TestInstallCommandHelpDescribesSupportedModes(t *testing.T) {
	var output bytes.Buffer
	installCmd.SetOut(&output)
	t.Cleanup(func() { installCmd.SetOut(nil) })
	if err := installCmd.Help(); err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"install [tool]", "--all", "dependency order", "retry-safe"} {
		if !bytes.Contains(output.Bytes(), []byte(text)) {
			t.Fatalf("help omits %q:\n%s", text, output.String())
		}
	}
}
