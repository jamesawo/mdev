package prerequisites

var registry []Prerequisite

// Register adds a prerequisite to the global registry.
func Register(p Prerequisite) {
	registry = append(registry, p)
}

// List returns all registered prerequisites.
func List() []Prerequisite {
	return registry
}
