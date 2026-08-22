package prerequisites

import (
	"context"

	"github.com/jamesawo/mdev/internal/ui/messages"
)

// Prerequisite defines a system requirement that must exist
// before mdev can install development tools.
type Prerequisite interface {

	// Name returns the display name used in CLI output.
	Name() string

	// Check verifies if the prerequisite is installed.
	Check() bool

	// Install installs the prerequisite if missing.
	Install() error
}

// State describes whether a system prerequisite is usable by mdev.
type State string

const (
	StateReady    State = messages.ReadinessStateReady
	StateMissing  State = messages.ReadinessStateMissing
	StateOutdated State = messages.ReadinessStateOutdated
	StateBroken   State = messages.ReadinessStateBroken
)

// StateChecker provides an error-aware, context-aware readiness check.
type StateChecker interface {
	Readiness(context.Context) (State, string, error)
}

// DependencyProvider declares prerequisite dependencies by canonical name.
type DependencyProvider interface {
	PrerequisiteDependencies() []string
}

// Remediator describes and performs an explicitly approved system change.
type Remediator interface {
	RemediationDescription() string
	RemediateContext(context.Context) error
}

// Verifier verifies a prerequisite after remediation.
type Verifier interface {
	VerifyContext(context.Context) error
}
