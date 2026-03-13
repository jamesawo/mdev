package drive

import "os"

func List() ([]string, error) {
	entries, err := os.ReadDir("/Volumes")
	if err != nil {
		return nil, err
	}

	var drives []string

	for _, entry := range entries {
		if entry.IsDir() {
			drives = append(drives, entry.Name())
		}
	}

	return drives, nil
}
