package uninstall

import "github.com/jamesawo/mdev/internal/tools"

// BuildPlan returns the uninstall plan in safe reverse dependency order.
func BuildPlan(target string) ([]string, error) {

	ordered, err := tools.ResolveOrder()
	if err != nil {
		return nil, err
	}

	affected := map[string]bool{target: true}
	for changed := true; changed; {
		changed = false
		for _, tool := range tools.List() {
			if affected[tool.Name()] {
				continue
			}
			for _, dependency := range tool.Dependencies() {
				if affected[dependency] {
					affected[tool.Name()] = true
					changed = true
					break
				}
			}
		}
	}

	// build reverse-safe uninstall order
	var plan []string

	for i := len(ordered) - 1; i >= 0; i-- {

		name := ordered[i].Name()
		if affected[name] {
			plan = append(plan, name)
		}
	}

	return plan, nil
}
