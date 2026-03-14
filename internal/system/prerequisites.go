package system

type Prerequisite interface {
	Name() string
	Check() bool
	Install() error
}
