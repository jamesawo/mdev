package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

const indent = "  "

var out io.Writer = os.Stdout
var mu sync.Mutex

func Section(title string) {
	mu.Lock()
	defer mu.Unlock()
	blankLocked()
	fmt.Fprintln(out, title)
}

func Success(name string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(out, "%s✓ %s\n", indent, name)
}

func Fail(name string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(out, "%s✗ %s\n", indent, name)
}

func Info(text string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(out, "%s%s\n", indent, text)
}

func Command(cmd string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(out, "%s%s\n", indent, cmd)
}

func Ask(text string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(out, "%s ", text)
}

func Blank() {
	mu.Lock()
	defer mu.Unlock()
	blankLocked()
}

func blankLocked() {
	fmt.Fprintln(out)
}

// ListItem prints a numbered list item.
// Example: "  1. scandisk"
func ListItem(index int, text string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(out, "%s%d. %s\n", indent, index, text)
}

func Indent(level int, text string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(out, "%s%s\n", strings.Repeat(indent, level), text)
}

func FormatIndent(level int, text string) string {
	return fmt.Sprintf("%s%s", strings.Repeat(indent, level), text)
}

func OverwriteLine(text string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(out, "\r\033[2K%s", text)
}

func ClearLine() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprint(out, "\r\033[2K")
}
