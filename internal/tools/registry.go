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
