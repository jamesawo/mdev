package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const indent = "  "

var out io.Writer = os.Stdout

func Section(title string) {
	Blank()
	fmt.Fprintln(out, title)
}

func Success(name string) {
	fmt.Fprintf(out, "%s✓ %s\n", indent, name)
}

func Fail(name string) {
	fmt.Fprintf(out, "%s✗ %s\n", indent, name)
}

func Info(text string) {
	fmt.Fprintf(out, "%s%s\n", indent, text)
}

func Command(cmd string) {
	fmt.Fprintf(out, "%s%s\n", indent, cmd)
}

func Ask(text string) {
	fmt.Fprintf(out, "%s ", text)
}

func Blank() {
	fmt.Fprintln(out)
}

// ListItem prints a numbered list item.
// Example: "  1. scandisk"
func ListItem(index int, text string) {
	fmt.Fprintf(out, "%s%d. %s\n", indent, index, text)
}

func Indent(level int, text string) {
	fmt.Fprintf(out, "%s%s\n", strings.Repeat(indent, level), text)
}
