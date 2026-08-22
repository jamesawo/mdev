package messages

import (
	"strings"
	"testing"
)

func TestSetupVolumeOptionIncludesNamePathAndReadOnlyState(t *testing.T) {
	writable := SetupVolumeOption("Work", "/Volumes/Work", true)
	if !strings.Contains(writable, "Work") || !strings.Contains(writable, "/Volumes/Work") || strings.Contains(writable, "read-only") {
		t.Fatalf("writable option = %q", writable)
	}
	readOnly := SetupVolumeOption("Archive", "/Volumes/Archive", false)
	if !strings.Contains(readOnly, "read-only") {
		t.Fatalf("read-only option = %q", readOnly)
	}
}

func TestSetupOutputFormatting(t *testing.T) {
	if got := SetupStorage("~/mdev"); got != "storage: ~/mdev" {
		t.Fatalf("SetupStorage() = %q", got)
	}
	if got := SetupSymlinkResolution("~/linked", "/Volumes/Data/mdev"); got != "~/linked  → /Volumes/Data/mdev" {
		t.Fatalf("SetupSymlinkResolution() = %q", got)
	}
	if SetupListCommand != "mdev list" {
		t.Fatalf("SetupListCommand = %q", SetupListCommand)
	}
}
