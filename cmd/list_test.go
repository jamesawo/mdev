package cmd

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	commandlist "github.com/jamesawo/mdev/internal/command/list"
)

func TestListCommandReturnsWorkflowErrorAfterWritingOutput(t *testing.T) {
	wantErr := errors.New("unknown status")
	runList = func(out io.Writer, _ commandlist.Options) error {
		_, _ = io.WriteString(out, "tools\n  ? example  unknown\n")
		return wantErr
	}
	t.Cleanup(func() { runList = defaultRunList })

	var output bytes.Buffer
	listCmd.SetOut(&output)
	t.Cleanup(func() { listCmd.SetOut(nil) })
	err := listCmd.RunE(listCmd, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("RunE() error = %v, want %v", err, wantErr)
	}
	if !strings.Contains(output.String(), "? example  unknown") {
		t.Fatalf("output = %q", output.String())
	}
}

func TestListCommandPassesJSONOption(t *testing.T) {
	listJSON = true
	t.Cleanup(func() { listJSON = false })
	var got commandlist.Options
	runList = func(_ io.Writer, options commandlist.Options) error {
		got = options
		return nil
	}
	t.Cleanup(func() { runList = defaultRunList })

	if err := listCmd.RunE(listCmd, nil); err != nil {
		t.Fatal(err)
	}
	if !got.JSON {
		t.Fatalf("RunE() options = %#v, want JSON", got)
	}
}

func TestListCommandHelpIsConciseAndReadOnly(t *testing.T) {
	var output bytes.Buffer
	listCmd.SetOut(&output)
	t.Cleanup(func() { listCmd.SetOut(nil) })
	if err := listCmd.Help(); err != nil {
		t.Fatal(err)
	}
	help := output.String()
	for _, text := range []string{"installation status", "system tools", "does not install"} {
		if !strings.Contains(help, text) {
			t.Fatalf("help omits %q:\n%s", text, help)
		}
	}
	if !strings.Contains(help, "--json") {
		t.Fatalf("list help omits --json:\n%s", help)
	}
}

func TestListCommandRejectsArguments(t *testing.T) {
	if err := listCmd.Args(listCmd, []string{"extra"}); err == nil {
		t.Fatal("list accepted an argument")
	}
}

func TestListCommandDoesNotDependOnYes(t *testing.T) {
	called := false
	runList = func(io.Writer, commandlist.Options) error {
		called = true
		return nil
	}
	t.Cleanup(func() { runList = defaultRunList })

	flag := rootCmd.PersistentFlags().Lookup("yes")
	if err := flag.Value.Set("true"); err != nil {
		t.Fatal(err)
	}
	flag.Changed = true
	t.Cleanup(func() {
		_ = flag.Value.Set("false")
		flag.Changed = false
	})
	if err := listCmd.RunE(listCmd, nil); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("list workflow was not called")
	}
}
