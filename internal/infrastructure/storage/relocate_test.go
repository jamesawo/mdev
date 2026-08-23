package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRelocateMovesStateAndCreatesManagedSymlink(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "home", "tool")
	target := filepath.Join(root, "managed", "tool")
	if err := os.MkdirAll(source, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "state"), []byte("kept"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Relocate(source, target); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(target, "state"))
	if err != nil || string(data) != "kept" {
		t.Fatalf("state = %q, error = %v", data, err)
	}
	got, err := filepath.EvalSymlinks(source)
	want, wantErr := filepath.EvalSymlinks(target)
	if err != nil || wantErr != nil || got != want {
		t.Fatalf("source resolves to %q, want %q, errors = %v, %v", got, want, err, wantErr)
	}
	if err := Relocate(source, target); err != nil {
		t.Fatalf("retry failed: %v", err)
	}
}

func TestRelocateRefusesTwoNonEmptyLocations(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	target := filepath.Join(root, "target")
	for _, path := range []string{source, target} {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(path, "state"), []byte(path), 0600); err != nil {
			t.Fatal(err)
		}
	}
	err := Relocate(source, target)
	if err == nil || !strings.Contains(err.Error(), "both") {
		t.Fatalf("error = %v", err)
	}
	for _, path := range []string{source, target} {
		if _, err := os.Stat(filepath.Join(path, "state")); err != nil {
			t.Fatalf("state at %s was changed: %v", path, err)
		}
	}
}

func TestRelocateRefusesUnexpectedSymlink(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	target := filepath.Join(root, "target")
	other := filepath.Join(root, "other")
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(other, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(other, source); err != nil {
		t.Fatal(err)
	}
	if err := Relocate(source, target); err == nil || !strings.Contains(err.Error(), "expected") {
		t.Fatalf("error = %v", err)
	}
}

func TestRelocateCreatesMissingSourceParents(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "home", ".local", "share", "tool", "state")
	target := filepath.Join(root, "managed", "tool")
	if err := Relocate(source, target); err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(source)
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatal(err)
	}
	if resolved != want {
		t.Fatalf("source resolves to %q, want %q", resolved, want)
	}
}

func TestRelocateRepairsExpectedDanglingSymlink(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "home", "machine")
	target := filepath.Join(root, "managed", "podman")
	if err := os.MkdirAll(filepath.Dir(source), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, source); err != nil {
		t.Fatal(err)
	}

	if err := Relocate(source, target); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		t.Fatalf("managed target was not recreated: %v", err)
	}
	resolved, err := filepath.EvalSymlinks(source)
	want, wantErr := filepath.EvalSymlinks(target)
	if err != nil || wantErr != nil || resolved != want {
		t.Fatalf("source resolves to %q, want %q, errors = %v, %v", resolved, want, err, wantErr)
	}
}
