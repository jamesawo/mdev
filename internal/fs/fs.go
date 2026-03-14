package fs

import (
	"os"
)

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsSymlink(path string) bool {

	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeSymlink != 0
}

func CreateSymlink(target, source string) error {
	return os.Symlink(target, source)
}

func Remove(path string) error {
	return os.Remove(path)
}

func Move(source, target string) error {
	return os.Rename(source, target)
}
