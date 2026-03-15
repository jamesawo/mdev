package install

import (
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
)

// RunAll installs every tool.
func RunAll(env *environment.Environment) error {

	ordered, err := tools.ResolveOrder()
	if err != nil {
		return err
	}

	return execute(env, ordered)
}

// RunSingle installs a specific tool and its dependencies.
func RunSingle(env *environment.Environment, name string) error {

	tool, ok := tools.Get(name)
	if !ok {
		return ErrUnknownTool(name)
	}

	plan, err := resolveSelection([]tools.Tool{tool})
	if err != nil {
		return err
	}

	return execute(env, plan)
}

// RunSelection installs tools chosen interactively.
func RunSelection(env *environment.Environment, names []string) error {

	var selected []tools.Tool

	for _, n := range names {

		t, ok := tools.Get(n)
		if !ok {
			return ErrUnknownTool(n)
		}

		selected = append(selected, t)
	}

	plan, err := resolveSelection(selected)
	if err != nil {
		return err
	}

	return execute(env, plan)
}
