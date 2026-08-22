package list

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/infrastructure/prerequisites"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

type status string

const (
	statusInstalled status = "installed"
	statusMissing   status = "missing"
	statusUnknown   status = "unknown"
)

type check struct {
	name   string
	verify func() (bool, error)
}

type result struct {
	name   string
	status status
	err    error
}

type dependencies struct {
	loadEnvironment func() (*environment.Environment, bool, error)
	stat            func(string) (os.FileInfo, error)
	systemChecks    func() []check
	toolChecks      func(*environment.Environment) []check
}

// UnknownStatusError reports that list completed its inventory but could not
// determine one or more statuses.
type UnknownStatusError struct {
	errors []error
}

func (e *UnknownStatusError) Error() string {
	return messages.ListUnknownStatuses
}

func (e *UnknownStatusError) Unwrap() []error {
	return e.errors
}

// Run prints all registered system prerequisites and tools with their current
// installation status. It observes configuration and machine state only.
func Run(out io.Writer) error {
	return run(out, productionDependencies())
}

func productionDependencies() dependencies {
	return dependencies{
		loadEnvironment: environment.Existing,
		stat:            os.Stat,
		systemChecks: func() []check {
			registered := prerequisites.List()
			checks := make([]check, 0, len(registered))
			for _, prerequisite := range registered {
				prerequisite := prerequisite
				checks = append(checks, check{
					name: prerequisite.Name(),
					verify: func() (bool, error) {
						return prerequisites.InstallationStatus(prerequisite)
					},
				})
			}
			return checks
		},
		toolChecks: func(env *environment.Environment) []check {
			registered := tools.List()
			checks := make([]check, 0, len(registered))
			for _, tool := range registered {
				tool := tool
				checks = append(checks, check{
					name: tool.Name(),
					verify: func() (bool, error) {
						return tools.InstallationStatus(tool, env)
					},
				})
			}
			return checks
		},
	}
}

func run(out io.Writer, deps dependencies) error {
	env, configured, err := deps.loadEnvironment()
	if err != nil {
		return err
	}
	if !configured {
		return errors.New(messages.ListNotConfigured)
	}

	info, err := deps.stat(env.StoragePath)
	if err != nil {
		return fmt.Errorf(messages.ListStorageUnavailable, env.StoragePath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf(messages.ListStorageNotDirectory, env.StoragePath)
	}

	systemResults := performChecks(deps.systemChecks())
	toolResults := performChecks(deps.toolChecks(env))
	render(out, systemResults, toolResults)

	var statusErrors []error
	for _, item := range append(systemResults, toolResults...) {
		if item.err != nil {
			statusErrors = append(statusErrors, fmt.Errorf("%s: %w", item.name, item.err))
		}
	}
	if len(statusErrors) > 0 {
		return &UnknownStatusError{errors: statusErrors}
	}
	return nil
}

func performChecks(checks []check) []result {
	sort.Slice(checks, func(i, j int) bool {
		left := strings.ToLower(checks[i].name)
		right := strings.ToLower(checks[j].name)
		if left == right {
			return checks[i].name < checks[j].name
		}
		return left < right
	})

	results := make([]result, 0, len(checks))
	for _, item := range checks {
		installed, err := item.verify()
		state := statusMissing
		if err != nil {
			state = statusUnknown
		} else if installed {
			state = statusInstalled
		}
		results = append(results, result{name: item.name, status: state, err: err})
	}
	return results
}

func render(out io.Writer, systemResults, toolResults []result) {
	printedSection := false
	if len(systemResults) > 0 {
		renderSection(out, messages.ListSystemTools, systemResults)
		printedSection = true
	}
	if len(toolResults) > 0 {
		if printedSection {
			fmt.Fprintln(out)
		}
		renderSection(out, messages.ListTools, toolResults)
		printedSection = true
	}

	unknown := append(unknownResults(systemResults), unknownResults(toolResults)...)
	if len(unknown) == 0 {
		return
	}
	if printedSection {
		fmt.Fprintln(out)
	}
	for _, item := range unknown {
		fmt.Fprintf(out, messages.ListUnknownDetail+"\n", item.name, item.err)
	}
}

func renderSection(out io.Writer, title string, results []result) {
	fmt.Fprintln(out, title)
	width := 0
	for _, item := range results {
		if len(item.name) > width {
			width = len(item.name)
		}
	}
	for _, item := range results {
		fmt.Fprintf(out, "  %s %-*s  %s\n", statusSymbol(item.status), width, item.name, item.status)
	}
}

func unknownResults(results []result) []result {
	unknown := make([]result, 0)
	for _, item := range results {
		if item.status == statusUnknown {
			unknown = append(unknown, item)
		}
	}
	return unknown
}

func statusSymbol(state status) string {
	switch state {
	case statusInstalled:
		return "✓"
	case statusMissing:
		return "○"
	default:
		return "?"
	}
}
