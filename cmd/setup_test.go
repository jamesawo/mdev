package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	commandsetup "github.com/jamesawo/mdev/internal/command/setup"
)

func TestSetupHelpOmitsYesAndDocumentsStoragePath(t *testing.T) {
	var output bytes.Buffer
	setupCmd.SetOut(&output)
	t.Cleanup(func() { setupCmd.SetOut(nil) })
	if err := setupCmd.Help(); err != nil {
		t.Fatal(err)
	}
	help := output.String()
	if strings.Contains(help, "--yes") {
		t.Fatalf("setup help includes unsupported --yes flag:\n%s", help)
	}
	if !strings.Contains(help, "--storage-path") {
		t.Fatalf("setup help omits --storage-path:\n%s", help)
	}
}

func TestSetupRejectsInheritedYes(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("yes")
	if err := flag.Value.Set("true"); err != nil {
		t.Fatal(err)
	}
	flag.Changed = true
	t.Cleanup(func() {
		_ = flag.Value.Set("false")
		flag.Changed = false
	})
	if err := setupCmd.PreRunE(setupCmd, nil); err == nil || !strings.Contains(err.Error(), "does not support") {
		t.Fatalf("PreRunE() error = %v", err)
	}
}

func TestSetupCommandDelegatesOptionsAndError(t *testing.T) {
	wantErr := errors.New("setup failed")
	var got commandsetup.Options
	runSetup = func(_ context.Context, _ commandsetup.Streams, options commandsetup.Options) error {
		got = options
		return wantErr
	}
	t.Cleanup(func() { runSetup = defaultRunSetup })
	setupStoragePath = "/tmp/example"
	t.Cleanup(func() { setupStoragePath = "" })

	err := setupCmd.RunE(setupCmd, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("RunE() error = %v, want %v", err, wantErr)
	}
	if got.StoragePath != setupStoragePath {
		t.Fatalf("RunE() options = %#v", got)
	}
}
