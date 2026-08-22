package tools

import "sort"

var registry = map[string]Tool{}

// Register adds a tool under its canonical name.
func Register(tool Tool) {
	registry[tool.Name()] = tool
}

// Get returns the tool registered under an exact canonical name.
func Get(name string) (Tool, bool) {
	tool, ok := registry[name]
	return tool, ok
}

// List returns all registered tools ordered by canonical name.
func List() []Tool {
	var registered []Tool

	for _, t := range registry {
		registered = append(registered, t)
	}
	sort.Slice(registered, func(i, j int) bool { return registered[i].Name() < registered[j].Name() })
	return registered
}

// ResolveSubset resolves named tools and their dependency closure in order.
func ResolveSubset(names []string) ([]Tool, error) {

	selected := make([]Tool, 0, len(names))
	for _, n := range names {
		tool, ok := Get(n)
		if !ok {
			return nil, ErrUnknownDependency(n)
		}
		selected = append(selected, tool)
	}
	return ResolveDependencies(selected)
}
