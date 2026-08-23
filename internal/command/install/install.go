package install

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/readiness"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

// Options describes the install mode selected by Cobra.
type Options struct {
	Tool      string
	All       bool
	AssumeYes bool
}

// Streams contains the command-owned input and output streams used by install.
type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type choice struct {
	label string
	name  string
}

// workflowDependencies groups external operations so orchestration remains testable.
type workflowDependencies struct {
	loadEnvironment func() (*environment.Environment, bool, error)
	validateStorage func(*environment.Environment) error
	listTools       func() []tools.Tool
	getTool         func(string) (tools.Tool, bool)
	resolve         func([]tools.Tool) ([]tools.Tool, error)
	status          func(tools.Tool, *environment.Environment) (bool, error)
	selectTools     func([]choice) ([]string, error)
	newReporter     func(io.Writer) progressReporter
	checkReadiness  func(context.Context, []string) ([]readiness.Item, error)
}

// Run validates, plans, confirms, and executes an install request.
func Run(ctx context.Context, streams Streams, options Options) error {
	return run(ctx, streams, options, defaultWorkflowDependencies())
}

// defaultWorkflowDependencies connects install to the real application boundaries.
func defaultWorkflowDependencies() workflowDependencies {
	return workflowDependencies{
		loadEnvironment: environment.Existing,
		validateStorage: environment.ValidateInstallStorage,
		listTools:       tools.List,
		getTool:         tools.Get,
		resolve:         tools.ResolveDependencies,
		status:          tools.InstallationStatus,
		newReporter:     func(out io.Writer) progressReporter { return newTextProgressReporter(out) },
		checkReadiness: func(ctx context.Context, names []string) ([]readiness.Item, error) {
			return readiness.CheckNames(ctx, names, nil)
		},
		selectTools: func(choices []choice) ([]string, error) {
			labels := make([]string, 0, len(choices))
			byLabel := make(map[string]string, len(choices))
			for _, item := range choices {
				labels = append(labels, item.label)
				byLabel[item.label] = item.name
			}
			selected, err := interactive.MultiSelect(messages.InstallSelectTools, labels)
			if err != nil {
				return nil, err
			}
			names := make([]string, 0, len(selected))
			for _, label := range selected {
				names = append(names, byLabel[label])
			}
			return names, nil
		},
	}
}

// run contains the install orchestration using supplied application boundaries.
func run(ctx context.Context, streams Streams, options Options, deps workflowDependencies) error {
	reporter := deps.newReporter(streams.Out)
	if cancelled(ctx) {
		return reporter.Cancelled()
	}
	env, configured, err := deps.loadEnvironment()
	if err != nil {
		return fmt.Errorf(messages.InstallConfigurationError, err)
	}
	if !configured {
		return errors.New(messages.InstallNotConfigured)
	}
	if err := deps.validateStorage(env); err != nil {
		return err
	}

	selected, err := selectRequested(ctx, env, options, deps)
	if err != nil {
		if cancelled(ctx) {
			return reporter.Cancelled()
		}
		return err
	}
	if len(selected) == 0 {
		return reporter.NoSelection()
	}
	plan, err := deps.resolve(selected)
	if err != nil {
		return fmt.Errorf(messages.InstallPlanError, err)
	}
	readinessItems, err := deps.checkReadiness(ctx, tools.SystemPrerequisites(plan))
	if err != nil {
		return fmt.Errorf(messages.InstallReadinessError, err)
	}
	for _, item := range readiness.Unready(readinessItems) {
		return fmt.Errorf(messages.InstallPrerequisiteMissing, item.Prerequisite.Name(), item.State)
	}
	states := make(map[string]bool, len(plan))
	pending := make([]tools.Tool, 0, len(plan))
	for _, tool := range plan {
		if cancelled(ctx) {
			return reporter.Cancelled()
		}
		installed, err := deps.status(tool, env)
		if err != nil {
			return fmt.Errorf(messages.InstallStatusError, tool.Name(), err)
		}
		states[tool.Name()] = installed
		if !installed {
			pending = append(pending, tool)
		}
	}
	if len(pending) > 0 {
		if err := reporter.Plan(pending); err != nil {
			return err
		}
		confirmer := confirmation.New(streams.In, streams.Out, options.AssumeYes)
		if !confirmer.AskDefaultNo(messages.InstallContinueQuestion) {
			return reporter.Cancelled()
		}
	}

	requested := make(map[string]bool, len(selected))
	for _, tool := range selected {
		requested[tool.Name()] = true
	}
	return execute(ctx, env, plan, states, requested, reporter)
}

// selectRequested resolves the explicit, all-tools, or interactive request mode.
func selectRequested(ctx context.Context, env *environment.Environment, options Options, deps workflowDependencies) ([]tools.Tool, error) {
	if options.All {
		return deps.listTools(), nil
	}
	if options.Tool != "" {
		tool, ok := deps.getTool(options.Tool)
		if !ok {
			return nil, fmt.Errorf(messages.InstallUnknownTool, options.Tool)
		}
		return []tools.Tool{tool}, nil
	}

	registered := deps.listTools()
	choices := make([]choice, 0, len(registered))
	for _, tool := range registered {
		if cancelled(ctx) {
			return nil, ctx.Err()
		}
		installed, err := deps.status(tool, env)
		if err != nil {
			return nil, fmt.Errorf(messages.InstallStatusError, tool.Name(), err)
		}
		label := tool.Name()
		if installed {
			label += messages.ToolsInstalledSuffix
		}
		choices = append(choices, choice{label: label, name: tool.Name()})
	}
	names, err := deps.selectTools(choices)
	if err != nil {
		return nil, err
	}
	selected := make([]tools.Tool, 0, len(names))
	for _, name := range names {
		tool, ok := deps.getTool(name)
		if !ok {
			return nil, fmt.Errorf(messages.InstallUnknownTool, name)
		}
		selected = append(selected, tool)
	}
	return selected, nil
}

// execute runs each incomplete tool through install, configure, and verify.
func execute(ctx context.Context, env *environment.Environment, plan []tools.Tool, states map[string]bool, requested map[string]bool, reporter progressReporter) error {
	for _, tool := range plan {
		if cancelled(ctx) {
			return reporter.Cancelled()
		}
		if states[tool.Name()] {
			if requested[tool.Name()] {
				if err := reporter.AlreadyInstalled(tool.Name(), true); err != nil {
					return err
				}
			}
			continue
		}
		for _, phase := range []struct {
			name   string
			action string
			run    func(context.Context, tools.Tool, *environment.Environment) error
		}{
			{name: messages.InstallPhaseInstall, action: messages.InstallActionInstall, run: tools.InstallContext},
			{name: messages.InstallPhaseConfigure, action: messages.InstallActionConfigure, run: tools.ConfigureContext},
			{name: messages.InstallPhaseVerify, action: messages.InstallActionVerify, run: tools.VerifyContext},
		} {
			if cancelled(ctx) {
				return reporter.Cancelled()
			}
			if err := reporter.Started(tool.Name(), phase.name); err != nil {
				return err
			}
			if err := phase.run(ctx, tool, env); err != nil {
				if cancelled(ctx) || errors.Is(err, context.Canceled) {
					return reporter.Cancelled()
				}
				return fmt.Errorf(messages.InstallPhaseError, phase.action, tool.Name(), err)
			}
		}
		if err := reporter.Completed(tool.Name()); err != nil {
			return err
		}
	}
	return nil
}

func cancelled(ctx context.Context) bool { return ctx.Err() != nil }
