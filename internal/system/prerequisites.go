package system

type Prerequisite struct {
	Name  string
	Check func() bool
}

func ListPrerequisites() []Prerequisite {

	return []Prerequisite{
		{
			Name: "brew",
			Check: func() bool {
				return CommandExists("brew")
			},
		},
		{
			Name: "curl",
			Check: func() bool {
				return CommandExists("curl")
			},
		},
	}
}
