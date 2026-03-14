package tools

import "fmt"

func ResolveOrder() ([]Tool, error) {

	visited := map[string]bool{}
	stack := map[string]bool{}
	order := []Tool{}

	var visit func(t Tool) error

	visit = func(t Tool) error {

		name := t.Name()

		if stack[name] {
			return fmt.Errorf("dependency cycle detected at %s", name)
		}

		if visited[name] {
			return nil
		}

		stack[name] = true

		for _, dep := range t.Dependencies() {

			depTool, ok := Get(dep)
			if !ok {
				return fmt.Errorf("unknown dependency %s", dep)
			}

			err := visit(depTool)
			if err != nil {
				return err
			}
		}

		stack[name] = false
		visited[name] = true

		order = append(order, t)

		return nil
	}

	for _, t := range List() {
		err := visit(t)
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}
