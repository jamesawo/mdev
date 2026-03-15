package interactive

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Select presents a numbered list of options to the user
// and waits for them to choose one.
//
// Behavior:
//   - The user selects an option by entering its number
//   - The function keeps prompting until valid input is received
//   - Entering "q" exits the selection
//
// Return values:
//
//	index >= 0 → selected option index
//	index = -1 → user chose to quit
func Select(options []string) (int, error) {

	reader := bufio.NewReader(os.Stdin)

	for {

		// Display options as a numbered list
		for i, option := range options {
			printer.ListItem(i+1, option)
		}

		// Display selection instructions
		printer.Info("Enter the number of your choice, or 'q' to quit:")
		printer.Ask(">")

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			return -1, err
		}

		// Normalize the input
		input = strings.TrimSpace(strings.ToLower(input))

		// Allow user to exit the selection
		if input == "q" {
			return -1, nil
		}

		// Convert the input string into an integer
		selection, err := strconv.Atoi(input)
		if err != nil {
			printer.Fail("Invalid input. Please enter a number.")
			continue
		}

		// Ensure the selected number is within the valid range
		if selection < 1 || selection > len(options) {
			printer.Fail("Selection out of range.")
			continue
		}

		// Convert the user-friendly number (1..n)
		// into a zero-based index used by Go slices
		return selection - 1, nil
	}
}
