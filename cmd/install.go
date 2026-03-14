package cmd

import (
	"fmt"

	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [tool]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Install a tool",
	Long: `
	Install a development tool into your local environment.

This command installs a supported tool and prepares it for use with the
current mdev environment configuration. The tool will be downloaded,
installed, and configured using the paths and settings defined in your
mdev environment.

Before running this command, your environment must be initialized and
validated using 'mdev doctor'. The install process depends on the
configured directories, tool paths, and system checks performed during
that step.

If the tool is already installed, the command will detect it and skip
the installation to avoid overwriting an existing setup.

Usage:
  mdev install [tool]

Arguments:
  tool    Name of the tool to install.

Behavior:
  • Validates that the environment is configured.
  • Checks whether the requested tool is supported.
  • Detects if the tool is already installed.
  • Runs the tool-specific installation process.

Examples:
  mdev install java
  mdev install gradle
  mdev install maven

Notes:
  Each tool provides its own installation logic. The command acts as a
  dispatcher that resolves the requested tool and executes its install
  routine using the current environment configuration.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		env, err := loadEnvironment()
		if err != nil {
			fmt.Println("Environment not configured. Run `mdev doctor` first.")
			return
		}

		if installAll {

			// install tools in the right order, based on their dependencies
			ordered, err := tools.ResolveOrder()
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, t := range ordered {

				err := installTool(env, t, map[string]bool{})
				if err != nil {
					fmt.Println(err)
				}
			}

			return
		}

		if len(args) == 0 {
			fmt.Println("Please specify a tool or use --all")
			return
		}

		name := args[0]
		tool, err := resolveTool(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = installTool(env, tool, map[string]bool{})
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
var installAll bool

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	installCmd.Flags().BoolVar(&installAll, "all", false, "Install all tools")
}

func resolveTool(name string) (tools.Tool, error) {

	tool, ok := tools.Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}

	return tool, nil
}

/*
*
Single Tool Install Flow
Run()

	└─ resolveTool()
	└─ installTool()
	     └─ IsInstalled() ✓ checked

with --all flag set
Run()

	└─ tools.List()
	    └─ installTool() --> checks dependencies
	         └─ IsInstalled() ✓ checked
*/
func installTool(env *environment.Environment, tool tools.Tool, visited map[string]bool) error {

	name := tool.Name()

	if visited[name] {
		return fmt.Errorf("dependency cycle detected at %s", name)
	}

	visited[name] = true

	for _, dep := range tool.Dependencies() {

		depTool, ok := tools.Get(dep)
		if !ok {
			return fmt.Errorf("missing dependency tool: %s", dep)
		}

		err := installTool(env, depTool, visited)
		if err != nil {
			return err
		}
	}

	if tool.IsInstalled(env) {
		fmt.Println("✓", name, "already installed")
		return nil
	}

	fmt.Println("Installing", name)

	err := tool.Install(env)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	err = tool.Configure(env)
	if err != nil {
		return fmt.Errorf("configuration failed: %w", err)
	}

	err = tool.Verify(env)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	fmt.Println("✓ Installed and configured", name)

	return nil
}
