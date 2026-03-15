package cmd

import (
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Show the dependency graph of supported tools",
	Long: `Graph displays the dependency relationships between development
tools managed by mdev.

Some tools depend on others to function correctly. For example,
Java requires SDKMAN to install and manage versions, while tools
like Maven and Gradle require Java to be present before they can
be installed.

This command prints the dependency graph so you can understand the
order in which tools will be installed when running:

  mdev install --all

The output shows each tool and the tools it depends on.

Example:

  mdev graph

Possible output:

  java -> sdkman
  maven -> java
  gradle -> java
  nvm
  podman
`,
	Run: func(cmd *cobra.Command, args []string) {

		printer.Section("Tool dependency graph")

		// Build reverse dependency graph:
		// dependency -> tools that depend on it
		graph := map[string][]string{}

		for _, t := range tools.List() {
			for _, dep := range t.Dependencies() {
				graph[dep] = append(graph[dep], t.Name())
			}
		}

		// Find root tools (tools with no dependencies)
		var roots []string
		for _, t := range tools.List() {
			if len(t.Dependencies()) == 0 {
				roots = append(roots, t.Name())
			}
		}

		// Print tree starting from each root
		for _, r := range roots {
			printTree(r, graph, 0)
		}
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
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
		printer.Info(indent + "└─ " + node)
	}

	for _, child := range graph[node] {
		printTree(child, graph, level+1)
	}
}
