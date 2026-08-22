package prerequisites

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestModernBashPathUsesVerifiedExecutableFromPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bash")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nprintf 5\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	got, err := ModernBashPath(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got != path {
		t.Fatalf("path = %q, want %q", got, path)
	}
}
