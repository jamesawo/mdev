package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AskYesNo ask a yes or no question like 'Install Homebrew now? (Y/n):'
func AskYesNo(question string) bool {

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s (Y/n): ", question)

	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(strings.ToLower(input))

	if input == "n" || input == "no" {
		return false
	}

	return true
}
