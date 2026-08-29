package uninstall

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/tools"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

// Options describes the uninstall request selected by Cobra.
type Options struct {
	Tool      string
	AssumeYes bool
}

// Streams contains command-owned input and output streams.
type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type workflowDependencies struct {
	loadEnvironment func() (*environment.Environment, error)
	buildPlan       func(string) ([]string, error)
	status          func(string, *environment.Environment) (bool, error)
	getTool         func(string) (tools.Tool, bool)
	stat            func(string) (os.FileInfo, error)
	removeAll       func(string) error
	uninstall       func(context.Context, tools.Tool, *environment.Environment) error
	newReporter     func(io.Writer) progressReporter
}

// Run resolves, confirms, and executes a managed uninstall request.
func Run(ctx context.Context, streams Streams, options Options) error {
	return run(ctx, streams, options, defaultWorkflowDependencies())
}

func defaultWorkflowDependencies() workflowDependencies {
	return workflowDependencies{
		loadEnvironment: environment.FromConfig,
		buildPlan:       BuildPlan,
		getTool:         tools.Get,
		stat:            os.Stat,
		removeAll:       os.RemoveAll,
		uninstall:       tools.UninstallContext,
		newReporter:     func(out io.Writer) progressReporter { return newTextProgressReporter(out) },
		status: func(name string, env *environment.Environment) (bool, error) {
			tool, ok := tools.Get(name)
			if !ok {
				return false, nil
			}
			return tools.InstallationStatus(tool, env)
		},
	}
}

func run(ctx context.Context, streams Streams, options Options, deps workflowDependencies) error {
	reporter := deps.newReporter(streams.Out)
	if ctx.Err() != nil {
		return reporter.Cancelled()
	}
	env, err := deps.loadEnvironment()
	if err != nil {
		return fmt.Errorf(messages.UninstallEnvironmentError, err)
	}
	resolved, err := deps.buildPlan(options.Tool)
	if err != nil {
		return fmt.Errorf(messages.UninstallPlanError, err)
	}
	plan, err := installedPlan(resolved, func(name string) (bool, error) {
		return deps.status(name, env)
	})
	if err != nil {
		return err
	}
	if len(plan) == 0 {
		return reporter.NotInstalled(options.Tool)
	}

	confirmer := confirmation.New(streams.In, streams.Out, options.AssumeYes)
	if len(plan) > 1 {
		if err := reporter.DependencyWarning(options.Tool, plan[:len(plan)-1]); err != nil {
			return err
		}
		if !confirmer.AskDefaultNo(messages.UninstallRemoveDependentsQuestion) {
			return reporter.Cancelled()
		}
	}
	if err := reporter.Plan(plan); err != nil {
		return err
	}

	directories, err := managedDirectories(plan, env, deps)
	if err != nil {
		return err
	}
	if err := reporter.Directories(directories); err != nil {
		return err
	}
	if !confirmer.AskDefaultNo(messages.UninstallContinueQuestion) {
		return reporter.Cancelled()
	}
	return execute(ctx, env, plan, reporter, deps)
}

func installedPlan(resolved []string, status func(string) (bool, error)) ([]string, error) {
	plan := make([]string, 0, len(resolved))
	for _, name := range resolved {
		installed, err := status(name)
		if err != nil {
			return nil, err
		}
		if installed {
			plan = append(plan, name)
		}
	}
	return plan, nil
}

func managedDirectories(plan []string, env *environment.Environment, deps workflowDependencies) ([]string, error) {
	var directories []string
	for _, name := range plan {
		tool, ok := deps.getTool(name)
		if !ok {
			continue
		}
		path := tool.StorageDir(env)
		if path == "" {
			continue
		}
		if _, err := deps.stat(path); err == nil {
			directories = append(directories, path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}
	return directories, nil
}
