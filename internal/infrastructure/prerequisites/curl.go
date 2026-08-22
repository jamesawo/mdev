package prerequisites

import (
	"errors"
	"os/exec"
)

type Curl struct{}

func (c *Curl) Name() string {
	return "curl"
}

func (c *Curl) Check() bool {
	installed, _ := c.InstallationStatus()
	return installed
}

func (c *Curl) InstallationStatus() (bool, error) {
	_, err := exec.LookPath("curl")
	if err == nil {
		return true, nil
	}
	if errors.Is(err, exec.ErrNotFound) {
		return false, nil
	}
	return false, err
}

func (c *Curl) Install() error {
	return nil
}

func init() {
	Register(&Curl{})
}
