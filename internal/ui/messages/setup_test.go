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
