package printer

import "fmt"

func Section(title string) {
	fmt.Println()
	fmt.Println(title)
}

func Success(name string) {
	fmt.Printf("  ✓ %s\n", name)
}

func Fail(name string) {
	fmt.Printf("  ✗ %s\n", name)
}

func Info(text string) {
	fmt.Printf("  %s\n", text)
}

func Command(cmd string) {
	fmt.Printf("  %s\n", cmd)
}

func Ask(text string) {
	fmt.Printf("%s ", text)
}

func Blank() {
	fmt.Println()
}

// ListItem prints a numbered list item.
// Example output: "  1. scandisk"
func ListItem(index int, text string) {
	fmt.Printf("  %d. %s\n", index, text)
}
