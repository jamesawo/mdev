package setup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/readiness"
	"github.com/jamesawo/mdev/internal/ui/confirmation"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Options contains setup inputs parsed by the CLI layer.
type Options struct {
	StoragePath string
}

// Streams contains setup's command-owned input and output streams.
type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type dependencies struct {
	setupInteractive   func() (*environment.Environment, error)
	existing           func() (*environment.Environment, bool, error)
	resolveStorage     func(string) (string, string, error)
	setupResolved      func(string) (*environment.Environment, error)
	validateResolved   func(string) error
	displayStoragePath func(string) string
	checkReadiness     func(context.Context, readiness.Reporter) ([]readiness.Item, error)
	remediate          func(context.Context, []readiness.Item, readiness.Reporter) error
	confirm            func(string) bool
	reporter           readiness.Reporter
}

type output interface {
	Section(string)
	Info(string)
	Blank()
	Command(string)
}

type terminalOutput struct{}

func (terminalOutput) Section(text string) { printer.Section(text) }
func (terminalOutput) Info(text string)    { printer.Info(text) }
func (terminalOutput) Blank()              { printer.Blank() }
func (terminalOutput) Command(text string) { printer.Command(text) }

// Run configures setup using process streams for compatibility with existing callers.
func Run(options Options) error {
	return RunContext(context.Background(), Streams{In: os.Stdin, Out: os.Stdout, Err: os.Stderr}, options)
}

// RunContext establishes system readiness and commits mdev storage configuration.
func RunContext(ctx context.Context, streams Streams, options Options) error {
	deps := productionDependencies()
	deps.confirm = confirmation.New(streams.In, streams.Out, false).AskDefaultNo
	deps.reporter = newReadinessReporter(streams.Out)
	return runContext(ctx, options, deps, terminalOutput{})
}

func productionDependencies() dependencies {
	return dependencies{
		setupInteractive:   setupInteractive,
		existing:           environment.Existing,
		resolveStorage:     environment.ResolveStoragePath,
		setupResolved:      environment.SetupResolved,
		validateResolved:   environment.ValidateResolvedSetupStorage,
		displayStoragePath: environment.DisplayPath,
		checkReadiness:     readiness.CheckAll,
		remediate:          readiness.Remediate,
	}
}

func run(options Options, deps dependencies, out output) error {
	return runContext(context.Background(), options, deps, out)
}

func runContext(ctx context.Context, options Options, deps dependencies, out output) error {
	if options.StoragePath == "" {
		return runInteractive(ctx, deps, out)
	}
	return runNonInteractive(ctx, options.StoragePath, deps, out)
}

func runInteractive(ctx context.Context, deps dependencies, out output) error {
	out.Section(messages.SetupTitle)
	env, err := deps.setupInteractive()
	if errors.Is(err, environment.ErrSetupCancelled) {
		out.Info(messages.SetupCancelled)
		return nil
	}
	if errors.Is(err, environment.ErrAlreadyConfigured) {
		if err := establishReadiness(ctx, deps, true); err != nil {
			return setupError(err)
		}
		printAlreadyConfigured(out, env, deps.displayStoragePath)
		return nil
	}
	if err != nil {
		return setupError(err)
	}
	if err := deps.validateResolved(env.StoragePath); err != nil {
		return setupError(err)
	}
	if err := establishReadiness(ctx, deps, true); err != nil {
		return setupError(err)
	}
	env, err = deps.setupResolved(env.StoragePath)
	if err != nil {
		return setupError(err)
	}
	printSuccess(out, env, deps.displayStoragePath)
	return nil
}

func runNonInteractive(ctx context.Context, storagePath string, deps dependencies, out output) error {
	existing, configured, err := deps.existing()
	if err != nil {
		return setupError(err)
	}
	if configured {
		return fmt.Errorf(messages.SetupAlreadyConfigured, existing.StoragePath)
	}

	resolved, _, err := deps.resolveStorage(storagePath)
	if err != nil {
		return setupError(err)
	}
	if err := deps.validateResolved(resolved); err != nil {
		return setupError(err)
	}
	if err := establishReadiness(ctx, deps, false); err != nil {
		return setupError(err)
	}
	env, err := deps.setupResolved(resolved)
	if err != nil {
		return setupError(err)
	}
	printSuccess(out, env, deps.displayStoragePath)
	return nil
}

func establishReadiness(ctx context.Context, deps dependencies, interactive bool) error {
	items, err := deps.checkReadiness(ctx, deps.reporter)
	if err != nil {
		return err
	}
	unready := readiness.Unready(items)
	if len(unready) == 0 {
		return nil
	}
	for _, item := range unready {
		if !item.Remediable() {
			return fmt.Errorf("system prerequisite %s is %s and requires manual remediation", item.Prerequisite.Name(), item.State)
		}
	}
	if !interactive {
		return errors.New(messages.SetupReadinessNonInteractive)
	}
	if deps.confirm == nil || !deps.confirm(messages.SetupReadinessApply) {
		return errors.New(messages.SetupReadinessDeclined)
	}
	return deps.remediate(ctx, items, deps.reporter)
}

func printAlreadyConfigured(out output, env *environment.Environment, displayPath func(string) string) {
	out.Info(messages.SetupComplete)
	out.Info(messages.SetupStorage(displayPath(env.StoragePath)))
	printNextStep(out)
}

func printSuccess(out output, env *environment.Environment, displayPath func(string) string) {
	out.Section(messages.SetupReady)
	out.Info(messages.SetupStorage(displayPath(env.StoragePath)))
	printNextStep(out)
}

func printNextStep(out output) {
	out.Blank()
	out.Info(messages.SetupNextStep)
	out.Command(messages.SetupListCommand)
}

func setupError(err error) error {
	return fmt.Errorf(messages.SetupError, messages.SetupFailed, err)
}
