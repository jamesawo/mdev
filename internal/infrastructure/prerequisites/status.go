package prerequisites

import "context"

// StatusChecker is implemented by prerequisites whose installation check can
// report an unexpected failure separately from a missing installation.
type StatusChecker interface {
	InstallationStatus() (bool, error)
}

// InstallationStatus preserves compatibility with existing Prerequisite
// implementations while allowing built-in prerequisites to expose failures.
func InstallationStatus(prerequisite Prerequisite) (bool, error) {
	checker, ok := prerequisite.(StatusChecker)
	if !ok {
		return prerequisite.Check(), nil
	}
	return checker.InstallationStatus()
}

// ReadinessStatus returns the richest readiness state supported by a prerequisite.
func ReadinessStatus(ctx context.Context, prerequisite Prerequisite) (State, string, error) {
	if checker, ok := prerequisite.(StateChecker); ok {
		return checker.Readiness(ctx)
	}
	installed, err := InstallationStatus(prerequisite)
	if err != nil {
		return StateBroken, "", err
	}
	if installed {
		return StateReady, "", nil
	}
	return StateMissing, "", nil
}
