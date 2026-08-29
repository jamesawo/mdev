package subprocess

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	maxCapturedBytes   = 64 * 1024
	maxDiagnosticBytes = 4 * 1024
	maxDiagnosticLines = 8
)

type managedOutputKey struct{}

// WithManagedOutput returns a child context that captures subprocess output
// for product-owned lifecycle presentation instead of streaming it directly.
func WithManagedOutput(ctx context.Context) context.Context {
	return context.WithValue(ctx, managedOutputKey{}, true)
}

// Run executes a cancellable subprocess. Normal callers inherit the process
// streams; managed callers capture successful chatter and retain bounded
// diagnostics on failure.
func Run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	if !managedOutput(ctx) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	stdout := &tailWriter{limit: maxCapturedBytes}
	stderr := &tailWriter{limit: maxCapturedBytes}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return withDiagnostic(err, diagnosticOutput(stdout.String(), stderr.String()))
	}
	return nil
}

// Check executes a cancellable subprocess without streaming successful output
// and returns bounded output only when the command fails.
func Check(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	stdout := &tailWriter{limit: maxCapturedBytes}
	stderr := &tailWriter{limit: maxCapturedBytes}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return withDiagnostic(err, diagnosticOutput(stdout.String(), stderr.String()))
	}
	return nil
}

func managedOutput(ctx context.Context) bool {
	managed, _ := ctx.Value(managedOutputKey{}).(bool)
	return managed
}

type commandError struct {
	err        error
	diagnostic string
}

func (e *commandError) Error() string {
	if e.diagnostic == "" {
		return e.err.Error()
	}
	return fmt.Sprintf("%v: %s", e.err, e.diagnostic)
}

func (e *commandError) Unwrap() error { return e.err }

func withDiagnostic(err error, output string) error {
	if err == nil {
		return nil
	}
	return &commandError{err: err, diagnostic: conciseDiagnostic(output)}
}

func diagnosticOutput(stdout, stderr string) string {
	if strings.TrimSpace(stderr) != "" {
		return stderr
	}
	return stdout
}

func conciseDiagnostic(output string) string {
	rawLines := strings.Split(strings.ReplaceAll(output, "\r", "\n"), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) > maxDiagnosticLines {
		lines = lines[len(lines)-maxDiagnosticLines:]
	}
	diagnostic := strings.TrimSpace(strings.Join(lines, "\n"))
	if len(diagnostic) > maxDiagnosticBytes {
		diagnostic = diagnostic[len(diagnostic)-maxDiagnosticBytes:]
		if index := strings.IndexByte(diagnostic, '\n'); index >= 0 {
			diagnostic = diagnostic[index+1:]
		}
	}
	return strings.TrimSpace(diagnostic)
}

type tailWriter struct {
	mu    sync.Mutex
	limit int
	data  []byte
}

func (w *tailWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.data = append(w.data, p...)
	if len(w.data) > w.limit {
		w.data = append([]byte(nil), w.data[len(w.data)-w.limit:]...)
	}
	return len(p), nil
}

func (w *tailWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return string(bytes.Clone(w.data))
}

var _ io.Writer = (*tailWriter)(nil)
var _ interface{ Unwrap() error } = (*commandError)(nil)
