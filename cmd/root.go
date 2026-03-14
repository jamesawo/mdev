package cmd

import (
	"os"

	// add tools via blank import here
	_ "github.com/jamesawo/mdev/internal/tools/gradle" //Gradle
	_ "github.com/jamesawo/mdev/internal/tools/java"   // JAVA
	_ "github.com/jamesawo/mdev/internal/tools/maven"  // Maven
	_ "github.com/jamesawo/mdev/internal/tools/nvm"    // NVM
	_ "github.com/jamesawo/mdev/internal/tools/podman" //Podman
	_ "github.com/jamesawo/mdev/internal/tools/sdkman" //SDKMAN

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mdev",
	Short: "Automate development environment setup on macOS",
	Long: `mdev is a command-line tool that automates the setup and management 
of a local development environment on macOS.

It installs development tools, configures them, and moves large caches
and data directories to an external drive to keep your system disk clean.

mdev manages tools such as Java, Maven, Gradle, Node, Podman and others,
handling installation, configuration, verification, and removal.

Typical workflow:

  1. Initialize your environment
     mdev doctor

  2. Install tools
     mdev install java
     mdev install maven

  3. Install the full development stack
     mdev install --all

mdev also provides commands to inspect your environment and visualize
tool dependencies.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mdev.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
