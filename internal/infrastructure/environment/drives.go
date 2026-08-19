package environment

import (
	"os"
	"path/filepath"
	"sort"
)

// ListExternalVolumes returns every mounted directory below the macOS volume
// root in deterministic, case-insensitive name order.
func ListExternalVolumes() ([]Volume, error) {
	volumeRoot := filepath.Join(string(filepath.Separator), "Volumes")
	return listVolumes(volumeRoot)
}

func listVolumes(volumeRoot string) ([]Volume, error) {
	entries, err := os.ReadDir(volumeRoot)
	if err != nil {
		return nil, err
	}
	volumes := make([]Volume, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil || !info.IsDir() {
			continue
		}
		path := filepath.Join(volumeRoot, entry.Name())
		volumes = append(volumes, Volume{Name: entry.Name(), Path: path, Writable: directoryWritable(path)})
	}
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Name < volumes[j].Name
	})
	return volumes, nil
}

func directoryWritable(path string) bool {
	probe, err := os.CreateTemp(path, ".mdev-volume-check-*")
	if err != nil {
		return false
	}
	name := probe.Name()
	if err := probe.Close(); err != nil {
		_ = os.Remove(name)
		return false
	}
	return os.Remove(name) == nil
}
