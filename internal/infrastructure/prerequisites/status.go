package prerequisites

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
