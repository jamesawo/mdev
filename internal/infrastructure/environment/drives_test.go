package environment

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListVolumesReturnsSortedWritableDirectories(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"Zulu", "Alpha"} {
		if err := os.Mkdir(filepath.Join(root, name), 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "not-a-volume"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	volumes, err := listVolumes(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(volumes) != 2 || volumes[0].Name != "Alpha" || volumes[1].Name != "Zulu" {
		t.Fatalf("volumes = %#v", volumes)
	}
	for _, volume := range volumes {
		if !volume.Writable {
			t.Fatalf("volume %q unexpectedly read-only", volume.Name)
		}
	}
}

func TestListVolumesSortsNamesCaseInsensitively(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"Zulu", "beta", "Alpha"} {
		if err := os.Mkdir(filepath.Join(root, name), 0755); err != nil {
			t.Fatal(err)
		}
	}
	volumes, err := listVolumes(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Alpha", "beta", "Zulu"}
	for index, name := range want {
		if volumes[index].Name != name {
			t.Fatalf("volumes = %#v", volumes)
		}
	}
}

func TestListVolumesReturnsEmptyList(t *testing.T) {
	volumes, err := listVolumes(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(volumes) != 0 {
		t.Fatalf("volumes = %#v", volumes)
	}
}

func TestListVolumesMarksReadOnlyVolume(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "Archive")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatal(err)
	}
	original := volumeWritable
	volumeWritable = func(candidate string) bool { return candidate != path }
	t.Cleanup(func() { volumeWritable = original })
	volumes, err := listVolumes(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(volumes) != 1 || volumes[0].Writable {
		t.Fatalf("volumes = %#v", volumes)
	}
}
