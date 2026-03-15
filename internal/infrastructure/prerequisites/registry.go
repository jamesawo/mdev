package prerequisites

// registry-based prerequisite system.
var prerequisites []Prerequisite

func Register(p Prerequisite) {
	prerequisites = append(prerequisites, p)
}

func List() []Prerequisite {
	return prerequisites
}
