package uninstall

import "github.com/jamesawo/mdev/internal/tools"

// BuildPlan returns the uninstall plan in safe reverse dependency order.
func BuildPlan(target string) ([]string, error) {

	ordered, err := tools.ResolveOrder()
	if err != nil {
		return nil, err
	}

	var affected []string

	// collect target and dependents
	for _, t := range tools.List() {

		if t.Name() == target {
			affected = append(affected, t.Name())
			continue
		}

		for _, dep := range t.Dependencies() {
			if dep == target {
				affected = append(affected, t.Name())
			}
		}
	}

	// ensure target included
	found := false
	for _, n := range affected {
		if n == target {
			found = true
		}
	}

	if !found {
		affected = append(affected, target)
	}

	// build reverse-safe uninstall order
	var plan []string

	for i := len(ordered) - 1; i >= 0; i-- {

		for _, a := range affected {
			if ordered[i].Name() == a {
				plan = append(plan, a)
			}
		}
	}

	return plan, nil
}
