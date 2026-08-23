package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestAppendTrailingBlankForHumanReadableOutput(t *testing.T) {
	var output bytes.Buffer
	command := &cobra.Command{Use: "example"}
	command.SetOut(&output)
	if err := appendTrailingBlank(command, nil); err != nil {
		t.Fatal(err)
	}
	if output.String() != "\n" {
		t.Fatalf("output = %q, want one trailing blank line", output.String())
	}
}

func TestAppendTrailingBlankPreservesJSONOutput(t *testing.T) {
	var output bytes.Buffer
	command := &cobra.Command{Use: "example"}
	command.SetOut(&output)
	command.Flags().Bool("json", false, "")
	if err := command.Flags().Set("json", "true"); err != nil {
		t.Fatal(err)
	}
	if err := appendTrailingBlank(command, nil); err != nil {
		t.Fatal(err)
	}
	if output.Len() != 0 {
		t.Fatalf("JSON output received trailing whitespace: %q", output.String())
	}
}
