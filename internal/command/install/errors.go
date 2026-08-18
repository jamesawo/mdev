package install

import (
	"errors"

	"github.com/jamesawo/mdev/internal/ui/messages"
)

func ErrUnknownTool(name string) error {
	return errors.New(messages.InstallUnknownTool(name))
}
