package cmd

import (
	"bytes"
	"errors"
	"go/parser"
	"go/token"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	commandversion "github.com/jamesawo/mdev/internal/command/version"
)

func TestVersionCommandDelegatesMetadataOptionsAndWriter(t *testing.T) {
	wantMetadata := commandversion.Metadata{Version: "1.2.3", Commit: "abc123", Built: "today"}
	versionMetadata = wantMetadata
	versionJSON = true
	t.Cleanup(func() {
		versionMetadata = commandversion.DefaultMetadata()
		versionJSON = false
	})

	var output bytes.Buffer
	versionCmd.SetOut(&output)
	t.Cleanup(func() { versionCmd.SetOut(nil) })

	var gotWriter io.Writer
	var gotMetadata commandversion.Metadata
	var gotOptions commandversion.Options
	runVersion = func(out io.Writer, metadata commandversion.Metadata, options commandversion.Options) error {
		gotWriter = out
		gotMetadata = metadata
		gotOptions = options
		return nil
	}
	t.Cleanup(func() { runVersion = defaultRunVersion })

	if err := versionCmd.RunE(versionCmd, nil); err != nil {
		t.Fatal(err)
	}
	if gotWriter != &output {
		t.Fatalf("writer = %T, want configured Cobra writer", gotWriter)
	}
	if gotMetadata != wantMetadata {
		t.Fatalf("metadata = %#v, want %#v", gotMetadata, wantMetadata)
	}
	if !gotOptions.JSON {
		t.Fatalf("options = %#v, want JSON", gotOptions)
	}
}

func TestVersionCommandPropagatesWorkflowError(t *testing.T) {
	wantErr := errors.New("write failed")
	runVersion = func(io.Writer, commandversion.Metadata, commandversion.Options) error {
		return wantErr
	}
	t.Cleanup(func() { runVersion = defaultRunVersion })

	if err := versionCmd.RunE(versionCmd, nil); !errors.Is(err, wantErr) {
		t.Fatalf("RunE() error = %v, want %v", err, wantErr)
	}
}

func TestVersionCommandHelpDocumentsJSONAndReadOnlyBehavior(t *testing.T) {
	var output bytes.Buffer
	versionCmd.SetOut(&output)
	t.Cleanup(func() { versionCmd.SetOut(nil) })
	if err := versionCmd.Help(); err != nil {
		t.Fatal(err)
	}
	help := output.String()
	for _, text := range []string{"--json", "read-only", "does not require mdev configuration"} {
		if !strings.Contains(help, text) {
			t.Fatalf("help omits %q:\n%s", text, help)
		}
	}
}

func TestVersionCommandRejectsArguments(t *testing.T) {
	if err := versionCmd.Args(versionCmd, []string{"extra"}); err == nil {
		t.Fatal("version accepted an argument")
	}
}

func TestRootVersionFlagUsesSharedConciseVersion(t *testing.T) {
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"--version"})
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
	})

	if _, err := rootCmd.ExecuteC(); err != nil {
		t.Fatal(err)
	}
	want := "mdev " + versionMetadata.Version + "\n\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestVersionOutputDoesNotIncludeAuthor(t *testing.T) {
	var output bytes.Buffer
	if err := commandversion.Run(&output, versionMetadata, commandversion.Options{}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(output.String()), "created by") || strings.Contains(output.String(), "James Aworo") {
		t.Fatalf("output includes author: %q", output.String())
	}
}

func TestVersionCommandRemainsThinCobraWiring(t *testing.T) {
	source, err := os.ReadFile("version.go")
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parser.ParseFile(token.NewFileSet(), "version.go", source, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}
	var imports []string
	for _, spec := range parsed.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			t.Fatal(err)
		}
		imports = append(imports, path)
	}
	wantImports := []string{
		"io",
		"github.com/jamesawo/mdev/internal/command/version",
		"github.com/jamesawo/mdev/internal/ui/messages",
		"github.com/spf13/cobra",
	}
	if !reflect.DeepEqual(imports, wantImports) {
		t.Fatalf("imports = %#v, want thin wiring imports %#v", imports, wantImports)
	}
	for _, forbidden := range []string{"fmt.", "json.", "os.", "exec.", "http."} {
		if strings.Contains(string(source), forbidden) {
			t.Fatalf("version.go contains business/runtime dependency %q", forbidden)
		}
	}
}
