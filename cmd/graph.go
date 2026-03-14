package cmd

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/tools"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Show tool dependency graph",
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
