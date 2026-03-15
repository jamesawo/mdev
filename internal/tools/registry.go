package tools

var registry = map[string]Tool{}

func Register(tool Tool) {
	registry[tool.Name()] = tool
}

func Get(name string) (Tool, bool) {
	tool, ok := registry[name]
	return tool, ok
}

func List() []Tool {
	var tools []Tool

	for _, t := range registry {
		tools = append(tools, t)
	}

	return tools
}

// ResolveSubset resolves dependency order for a subset of tools.
// todo: should this be here?
func ResolveSubset(names []string) ([]Tool, error) {

	ordered, err := ResolveOrder()
	if err != nil {
		return nil, err
	}

	selected := map[string]bool{}
	for _, n := range names {
		selected[n] = true
	}

	var result []Tool

	for _, t := range ordered {
		if selected[t.Name()] {
			result = append(result, t)
		}
	}

	return result, nil
}
