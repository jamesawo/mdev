package cmd

import (
	"fmt"
	"strconv"

	"github.com/jamesawo/mdev/internal/drive"
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

		fmt.Println("Selected drive:", selected)
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
