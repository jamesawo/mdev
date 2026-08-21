package environment

import "errors"

var (
	ErrSetupCancelled    = errors.New("setup cancelled")
	ErrAlreadyConfigured = errors.New("mdev is already configured")
)
