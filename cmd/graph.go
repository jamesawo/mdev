package cmd

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/tools"
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

This command is useful for debugging tool installation order and
understanding how mdev resolves dependencies during automated setup.`,
	Run: func(cmd *cobra.Command, args []string) {

		for _, t := range tools.List() {

			deps := t.Dependencies()

			if len(deps) == 0 {
				fmt.Println(t.Name())
				continue
			}

			for _, d := range deps {
				fmt.Printf("%s -> %s\n", t.Name(), d)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
