package environment

import (
	"os"
	"path/filepath"
)

// listExternalVolumes returns mounted directories from the macOS volume root.
func listExternalVolumes() ([]string, error) {
	volumeRoot := filepath.Join(string(filepath.Separator), "Volumes")

	entries, err := os.ReadDir(volumeRoot)
	if err != nil {
		return nil, err
	}

	var volumes []string

	for _, e := range entries {
		if e.IsDir() {
			volumes = append(volumes, filepath.Join(volumeRoot, e.Name()))
		}
	}

	return volumes, nil
}
