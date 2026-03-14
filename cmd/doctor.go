package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jamesawo/mdev/internal/config"
	"github.com/jamesawo/mdev/internal/drive"
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/jamesawo/mdev/internal/system"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Validate and initialize the mdev environment",
	Long: `Doctor validates that your development environment is correctly
configured for use with mdev.

It performs several checks including:

  • System prerequisites (brew, curl, etc.)
  • Existing mdev environment configuration
  • External storage availability
  • Tool installation status

If no environment has been configured yet, doctor will guide you
through the setup process and allow you to choose an external drive
where development tool data and caches will be stored.

Typical usage:

  mdev doctor

This is usually the first command you run on a new machine before
installing any development tools.`,
	Run: func(cmd *cobra.Command, args []string) {

		checkSystemPrerequisites()

		if checkExistingEnvironment() {
			return
		}

		setupEnvironment()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doctorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doctorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkExistingEnvironment() bool {

	env, err := loadEnvironment()
	if err != nil {
		return false
	}

	if !ensureDataRoot(env) {
		return true
	}

	printEnvironmentStatus(env)
	checkTools(env)

	return true
}

func setupEnvironment() {

	drives, err := drive.List()
	if err != nil {
		fmt.Println("Error reading drives:", err)
		return
	}

	fmt.Println("Available drives:")

	for i, d := range drives {
		fmt.Printf("%d. %s\n", i+1, d)
	}

	fmt.Print("Select a drive: ")

	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return
	}

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(drives) {
		fmt.Println("Invalid selection")
		return
	}

	selected := drives[index-1]

	path := "/Volumes/" + selected

	err = config.SaveExternalDrive(path)
	if err != nil {
		fmt.Println("Failed to save configuration:", err)
		return
	}

	fmt.Println("Configuration saved.")
	fmt.Println("External drive:", path)

	env := environment.New(path)

	err = environment.CreateDataRoot(env)
	if err != nil {
		fmt.Println("Failed to create data directory:", err)
	}
}

func checkTools(env *environment.Environment) {

	fmt.Println()
	fmt.Println("Tools:")

	for _, t := range tools.List() {

		if t.IsInstalled(env) {
			fmt.Println("✓", t.Name(), "installed")
		} else {
			fmt.Println("✗", t.Name(), "not installed")
		}
	}
}

func loadEnvironment() (*environment.Environment, error) {
	return environment.FromConfig()
}

func ensureDataRoot(env *environment.Environment) bool {

	err := environment.CreateDataRoot(env)
	if err != nil {
		fmt.Println("Failed to ensure data directory:", err)
		return false
	}

	return true
}

func printEnvironmentStatus(env *environment.Environment) {

	fmt.Println("Environment status:")
	fmt.Println("✓ External drive:", env.ExternalDrive)

	_, err := os.Stat(env.DataRoot)
	if err == nil {
		fmt.Println("✓ Data directory:", env.DataRoot)
	} else {
		fmt.Println("✗ Data directory missing:", env.DataRoot)
	}
}

func checkSystemPrerequisites() {

	fmt.Println("System prerequisites:")

	for _, p := range system.List() {

		if p.Check() {
			fmt.Println("✓", p.Name(), "installed")
			continue
		}

		fmt.Println("✗", p.Name(), "missing")

		err := p.Install()
		if err != nil {
			fmt.Println("failed to install", p.Name(), err)
		}
	}

	fmt.Println()
}
