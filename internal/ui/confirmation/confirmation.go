package confirmation

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Confirmer struct {
	assumeYes bool
	input     *bufio.Reader
	output    io.Writer
}

var defaultConfirmer = New(os.Stdin, os.Stdout, false)

func New(input io.Reader, output io.Writer, assumeYes bool) *Confirmer {
	return &Confirmer{
		assumeYes: assumeYes,
		input:     bufio.NewReader(input),
		output:    output,
	}
}

func Configure(assumeYes bool) {
	defaultConfirmer.assumeYes = assumeYes
}

func Ask(question string) bool {
	return defaultConfirmer.Ask(question)
}

func (c *Confirmer) Ask(question string) bool {
	if c.assumeYes {
		return true
	}

	fmt.Fprintf(c.output, "%s (Y/n): ", question)

	input, err := c.input.ReadString('\n')
	if err != nil && input == "" {
		return false
	}
	answer := strings.TrimSpace(strings.ToLower(input))

	return answer != "n" && answer != "no"
}
