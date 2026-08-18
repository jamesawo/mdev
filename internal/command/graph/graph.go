package graph

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

func Run() {

	printer.Section(messages.GraphTitle)

	graph := map[string][]string{}

	for _, t := range tools.List() {
		for _, dep := range t.Dependencies() {
			graph[dep] = append(graph[dep], t.Name())
		}
	}

	var roots []string

	for _, t := range tools.List() {
		if len(t.Dependencies()) == 0 {
			roots = append(roots, t.Name())
		}
	}

	for _, r := range roots {
		printTree(r, graph, 0)
	}
}

// printTree recursively prints the dependency tree.
// node: current tool
// graph: reverse dependency map
// level: indentation level
func printTree(node string, graph map[string][]string, level int) {

	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}

	if level == 0 {
		printer.Info(node)
	} else {
		printer.Info(fmt.Sprintf("%s %s %s", indent, messages.GraphIndentMark, node))
	}

	for _, child := range graph[node] {
		printTree(child, graph, level+1)
	}
}
