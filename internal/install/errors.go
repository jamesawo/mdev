package install

import "fmt"

func ErrUnknownTool(name string) error {
	return fmt.Errorf("unknown tool: %s", name)
}
