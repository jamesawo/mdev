package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jamesawo/mdev/internal/ui/printer"
)

// RadioSelect presents options as radio buttons and lets the user
// move the selection using j/k or arrow keys and confirm with Enter.
//
// Controls:
//
//	j / ↓  → move down
//	k / ↑  → move up
//	enter  → select
//	q      → quit
//
// Returns:
//
//	index >= 0 → selected option
//	index = -1 → user quit
func RadioSelect(title string, options []string) (int, error) {

	reader := bufio.NewReader(os.Stdin)

	selected := 0

	for {

		printer.Section(title)

		for i, option := range options {

			cursor := "( )"

			if i == selected {
				cursor = "(*)"
			}

			printer.Info(fmt.Sprintf("%s %s", cursor, option))
		}

		printer.Blank()
		printer.Info("Use arrows keys to move, Enter to select, q to quit")
		printer.Ask(">")

		input, err := reader.ReadString('\n')
		if err != nil {
			return -1, err
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {

		case "q":
			return -1, nil

		case "j":
			if selected < len(options)-1 {
				selected++
			}

		case "k":
			if selected > 0 {
				selected--
			}

		case "":
			return selected, nil
		}

		fmt.Print("\033[H\033[2J") // clear screen
	}
}
