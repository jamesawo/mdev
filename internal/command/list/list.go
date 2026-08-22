package list

import (
	"encoding/json"
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

// Options controls how list results are encoded.
type Options struct {
	JSON bool
}

type check struct {
	name   string
	verify func() (bool, error)
}

type result struct {
	name   string
	status status
	err    error
}

type jsonEntry struct {
	Name   string `json:"name"`
	Status status `json:"status"`
	Error  string `json:"error,omitempty"`
}

type jsonDocument struct {
	SystemTools []jsonEntry `json:"system_tools"`
	Tools       []jsonEntry `json:"tools"`
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
func Run(out io.Writer, options Options) error {
	return run(out, options, productionDependencies())
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

func run(out io.Writer, options Options, deps dependencies) error {
	env, err := loadEnvironment(deps)
	if err != nil {
		return err
	}

	systemChecks := sortedChecks(deps.systemChecks())
	toolChecks := sortedChecks(deps.toolChecks(env))
	if options.JSON {
		return runJSON(out, systemChecks, toolChecks)
	}
	return runHuman(out, systemChecks, toolChecks)
}

func loadEnvironment(deps dependencies) (*environment.Environment, error) {
	env, configured, err := deps.loadEnvironment()
	if err != nil {
		return nil, err
	}
	if !configured {
		return nil, errors.New(messages.ListNotConfigured)
	}

	info, err := deps.stat(env.StoragePath)
	if err != nil {
		return nil, fmt.Errorf(messages.ListStorageUnavailable, env.StoragePath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf(messages.ListStorageNotDirectory, env.StoragePath)
	}
	return env, nil
}

func runHuman(out io.Writer, systemChecks, toolChecks []check) error {
	var allResults []result
	printedSection := false

	for _, section := range []struct {
		title  string
		checks []check
	}{
		{title: messages.ListSystemTools, checks: systemChecks},
		{title: messages.ListTools, checks: toolChecks},
	} {
		if len(section.checks) == 0 {
			continue
		}
		if printedSection {
			if _, err := fmt.Fprintln(out); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(out, section.title); err != nil {
			return err
		}
		printedSection = true
		width := checkNameWidth(section.checks)
		for _, item := range section.checks {
			itemResult := performCheck(item)
			allResults = append(allResults, itemResult)
			if _, err := fmt.Fprintf(out, "  %s %-*s  %s\n", statusSymbol(itemResult.status), width, itemResult.name, itemResult.status); err != nil {
				return err
			}
		}
	}

	unknown := unknownResults(allResults)
	if len(unknown) > 0 && printedSection {
		if _, err := fmt.Fprintln(out); err != nil {
			return err
		}
	}
	for _, item := range unknown {
		if _, err := fmt.Fprintf(out, messages.ListUnknownDetail+"\n", item.name, conciseError(item.err)); err != nil {
			return err
		}
	}
	return unknownStatusError(unknown)
}

func runJSON(out io.Writer, systemChecks, toolChecks []check) error {
	systemResults := performChecks(systemChecks)
	toolResults := performChecks(toolChecks)
	document := jsonDocument{
		SystemTools: jsonEntries(systemResults),
		Tools:       jsonEntries(toolResults),
	}
	encoded, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("encode list JSON: %w", err)
	}
	encoded = append(encoded, '\n')
	if _, err := out.Write(encoded); err != nil {
		return err
	}
	return unknownStatusError(append(unknownResults(systemResults), unknownResults(toolResults)...))
}

func sortedChecks(checks []check) []check {
	sort.Slice(checks, func(i, j int) bool {
		left := strings.ToLower(checks[i].name)
		right := strings.ToLower(checks[j].name)
		if left == right {
			return checks[i].name < checks[j].name
		}
		return left < right
	})
	return checks
}

func performChecks(checks []check) []result {
	results := make([]result, 0, len(checks))
	for _, item := range checks {
		results = append(results, performCheck(item))
	}
	return results
}

func performCheck(item check) result {
	installed, err := item.verify()
	state := statusMissing
	if err != nil {
		state = statusUnknown
	} else if installed {
		state = statusInstalled
	}
	return result{name: item.name, status: state, err: err}
}

func jsonEntries(results []result) []jsonEntry {
	entries := make([]jsonEntry, 0, len(results))
	for _, item := range results {
		entry := jsonEntry{Name: item.name, Status: item.status}
		if item.err != nil {
			entry.Error = conciseError(item.err)
		}
		entries = append(entries, entry)
	}
	return entries
}

func unknownStatusError(unknown []result) error {
	if len(unknown) == 0 {
		return nil
	}
	statusErrors := make([]error, 0, len(unknown))
	for _, item := range unknown {
		statusErrors = append(statusErrors, fmt.Errorf("%s: %w", item.name, item.err))
	}
	return &UnknownStatusError{errors: statusErrors}
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

func checkNameWidth(checks []check) int {
	width := 0
	for _, item := range checks {
		if len(item.name) > width {
			width = len(item.name)
		}
	}
	return width
}

func conciseError(err error) string {
	detail := strings.TrimSpace(err.Error())
	if newline := strings.IndexByte(detail, '\n'); newline >= 0 {
		detail = detail[:newline]
	}
	const maxLength = 240
	if len(detail) > maxLength {
		detail = detail[:maxLength-3] + "..."
	}
	return detail
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
