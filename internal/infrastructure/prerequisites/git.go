package prerequisites

import (
	"errors"
	"os/exec"
)

type Git struct{}

func (Git) Name() string { return "git" }
func (Git) Check() bool {
	installed, _ := (Git{}).InstallationStatus()
	return installed
}
func (Git) InstallationStatus() (bool, error) {
	_, err := exec.LookPath("git")
	if err == nil {
		return true, nil
	}
	if errors.Is(err, exec.ErrNotFound) {
		return false, nil
	}
	return false, err
}
func (Git) Install() error { return nil }

func init() {
	Register(Git{})
}
