package environment

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateInstallStorageRequiresExistingDirectory(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	err := ValidateInstallStorage(New(missing))
	if err == nil || !strings.Contains(err.Error(), missing) {
		t.Fatalf("error = %v", err)
	}
	if _, statErr := os.Stat(missing); !os.IsNotExist(statErr) {
		t.Fatalf("storage was created: %v", statErr)
	}
}

func TestValidateInstallStorageAcceptsWritableDirectoryWithoutResidue(t *testing.T) {
	root := t.TempDir()
	if err := ValidateInstallStorage(New(root)); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("probe residue = %v", entries)
	}
}

func TestValidateInstallStorageRejectsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "storage")
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := ValidateInstallStorage(New(path)); err == nil || !strings.Contains(err.Error(), "directory") {
		t.Fatalf("error = %v", err)
	}
}
