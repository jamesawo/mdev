package environment

import "os"

// listExternalDrives returns directories under /Volumes.
func listExternalDrives() ([]string, error) {

	entries, err := os.ReadDir("/Volumes")
	if err != nil {
		return nil, err
	}

	var drives []string

	for _, e := range entries {
		if e.IsDir() {
			drives = append(drives, e.Name())
		}
	}

	return drives, nil
}
