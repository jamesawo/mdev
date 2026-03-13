package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jamesawo/mdev/internal/config"
	"github.com/jamesawo/mdev/internal/drive"
	"github.com/jamesawo/mdev/internal/environment"
	"github.com/spf13/cobra"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

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

	env, err := environment.FromConfig()
	if err != nil {
		return false
	}

	err = environment.CreateDataRoot(env)
	if err != nil {
		fmt.Println("Failed to ensure data directory:", err)
		return true
	}

	fmt.Println("Environment status:")
	fmt.Println("✓ External drive:", env.ExternalDrive)

	_, err = os.Stat(env.DataRoot)
	if err == nil {
		fmt.Println("✓ Data directory:", env.DataRoot)
	} else {
		fmt.Println("✗ Data directory missing:", env.DataRoot)
	}

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
