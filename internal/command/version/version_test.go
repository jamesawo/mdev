package version

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/parser"
	"go/token"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestRunTextUsesDeterministicDefaults(t *testing.T) {
	const want = "mdev 0.2.0\ncommit: unknown\nbuilt: unknown\n"

	for _, metadata := range []Metadata{{}, DefaultMetadata()} {
		var first bytes.Buffer
		if err := Run(&first, metadata, Options{}); err != nil {
			t.Fatal(err)
		}
		var second bytes.Buffer
		if err := Run(&second, metadata, Options{}); err != nil {
			t.Fatal(err)
		}
		if first.String() != want {
			t.Fatalf("output = %q, want %q", first.String(), want)
		}
		if second.String() != first.String() {
			t.Fatalf("repeated output changed: %q then %q", first.String(), second.String())
		}
	}
}

func TestRunTextUsesSuppliedMetadata(t *testing.T) {
	metadata := Metadata{Version: "1.2.3", Commit: "a81f37c", Built: "2026-08-22"}
	var output bytes.Buffer
	if err := Run(&output, metadata, Options{}); err != nil {
		t.Fatal(err)
	}
	const want = "mdev 1.2.3\ncommit: a81f37c\nbuilt: 2026-08-22\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestRunJSONUsesSameMetadataWithoutHumanText(t *testing.T) {
	metadata := Metadata{Version: "1.2.3", Commit: "a81f37c", Built: "2026-08-22"}
	var output bytes.Buffer
	if err := Run(&output, metadata, Options{JSON: true}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(output.String(), "\n") {
		t.Fatalf("JSON lacks trailing newline: %q", output.String())
	}

	var document map[string]any
	if err := json.Unmarshal(output.Bytes(), &document); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, output.String())
	}
	want := map[string]any{
		"version": metadata.Version,
		"commit":  metadata.Commit,
		"built":   metadata.Built,
	}
	if !reflect.DeepEqual(document, want) {
		t.Fatalf("document = %#v, want %#v", document, want)
	}
	for _, humanText := range []string{"mdev ", "commit:", "built:"} {
		if strings.Contains(output.String(), humanText) {
			t.Fatalf("JSON contains human text %q: %s", humanText, output.String())
		}
	}
}

func TestRunAppliesFallbacksConsistentlyAcrossModes(t *testing.T) {
	metadata := Metadata{Version: "2.0.0"}

	var textOutput bytes.Buffer
	if err := Run(&textOutput, metadata, Options{}); err != nil {
		t.Fatal(err)
	}
	if textOutput.String() != "mdev 2.0.0\ncommit: unknown\nbuilt: unknown\n" {
		t.Fatalf("text output = %q", textOutput.String())
	}

	var jsonOutput bytes.Buffer
	if err := Run(&jsonOutput, metadata, Options{JSON: true}); err != nil {
		t.Fatal(err)
	}
	var document Metadata
	if err := json.Unmarshal(jsonOutput.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	want := Metadata{Version: "2.0.0", Commit: "unknown", Built: "unknown"}
	if document != want {
		t.Fatalf("JSON metadata = %#v, want %#v", document, want)
	}
}

func TestRunReturnsWriterFailure(t *testing.T) {
	wantErr := errors.New("writer unavailable")
	err := Run(errorWriter{err: wantErr}, DefaultMetadata(), Options{JSON: true})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want wrapped %v", err, wantErr)
	}
}

func TestRunReturnsShortWriterFailure(t *testing.T) {
	err := Run(shortWriter{}, DefaultMetadata(), Options{})
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("Run() error = %v, want %v", err, io.ErrShortWrite)
	}
}

func TestWorkflowHasNoRuntimeEnvironmentDependencies(t *testing.T) {
	parsed, err := parser.ParseFile(token.NewFileSet(), "version.go", nil, parser.ImportsOnly)
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
	want := []string{
		"encoding/json",
		"fmt",
		"io",
		"github.com/jamesawo/mdev/internal/ui/messages",
	}
	if !reflect.DeepEqual(imports, want) {
		t.Fatalf("imports = %#v, want dependency-free workflow imports %#v", imports, want)
	}
}

type errorWriter struct {
	err error
}

type shortWriter struct{}

func (shortWriter) Write(contents []byte) (int, error) {
	return len(contents) - 1, nil
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}
