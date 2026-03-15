package prerequisites

type Git struct{}

func (Git) Name() string   { return "git" }
func (Git) Check() bool    { return true }
func (Git) Install() error { return nil }

func init() {
	Register(Git{})
}
