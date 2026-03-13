package tools

import "github.com/jamesawo/mdev/internal/environment"

type Tool interface {
	Name() string
	Description() string

	IsInstalled(env *environment.Environment) bool

	Install(env *environment.Environment) error
	Configure(env *environment.Environment) error
	Verify(env *environment.Environment) error

	Uninstall(env *environment.Environment) error
}
