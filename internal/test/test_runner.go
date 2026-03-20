package exec

import "context"

// FakeRunner is used in tests so we don't run real system commands.
type FakeRunner struct {
	Calls []string
}

// Run records the command instead of executing it.
func (f *FakeRunner) Run(ctx context.Context, name string, args ...string) error {

	call := name
	for _, a := range args {
		call += " " + a
	}

	f.Calls = append(f.Calls, call)

	return nil
}
