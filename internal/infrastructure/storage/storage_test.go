package storage

import (
	"path/filepath"
	"testing"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

func TestToolStorageIsDirectlyBelowStoragePath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "mdev")
	env := environment.New(root)

	got := ToolDir(env, "maven")
	want := filepath.Join(root, "maven")
	if got != want {
		t.Fatalf("ToolDir() = %q, want %q", got, want)
	}
	if got == filepath.Join(root, "data", "maven") {
		t.Fatal("ToolDir() retained obsolete data directory")
	}
}
