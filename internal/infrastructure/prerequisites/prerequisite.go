package prerequisites

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
