package messages

import "fmt"

const (
	UpCmdShortDescription = "Start development services (like podman, ollama)"
	UpCmdLongDescription  = `Start the development environment.

mdev up starts the runtime services and tools required during development.

This typically includes:
- podman machine
- ollama service


Use this command at the beginning of your development session.`
)

func UpServiceFailedToStart(name string, err error) error {
	return fmt.Errorf("%s failed to start: %w", name, err)
}

func UpServiceFailedToStop(name string, err error) error {
	return fmt.Errorf("%s failed to stop: %w", name, err)
}
