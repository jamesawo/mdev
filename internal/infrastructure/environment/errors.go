package environment

import "errors"

import "github.com/jamesawo/mdev/internal/ui/messages"

var (
	ErrSetupCancelled    = errors.New(messages.SetupCancelledError)
	ErrAlreadyConfigured = errors.New(messages.SetupConfiguredError)
)
