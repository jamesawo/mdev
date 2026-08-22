package setup

import (
	"errors"
	"fmt"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

// Options contains setup inputs parsed by the CLI layer.
type Options struct {
	StoragePath string
}

type dependencies struct {
	setupInteractive   func() (*environment.Environment, error)
	existing           func() (*environment.Environment, bool, error)
	resolveStorage     func(string) (string, string, error)
	setupResolved      func(string) (*environment.Environment, error)
	displayStoragePath func(string) string
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

// Run configures mdev storage or reports the existing configuration.
func Run(options Options) error {
	return run(options, productionDependencies(), terminalOutput{})
}

func productionDependencies() dependencies {
	return dependencies{
		setupInteractive:   setupInteractive,
		existing:           environment.Existing,
		resolveStorage:     environment.ResolveStoragePath,
		setupResolved:      environment.SetupResolved,
		displayStoragePath: environment.DisplayPath,
	}
}

func run(options Options, deps dependencies, out output) error {
	if options.StoragePath == "" {
		return runInteractive(deps, out)
	}
	return runNonInteractive(options.StoragePath, deps, out)
}

func runInteractive(deps dependencies, out output) error {
	out.Section(messages.SetupTitle)
	env, err := deps.setupInteractive()
	if errors.Is(err, environment.ErrSetupCancelled) {
		out.Info(messages.SetupCancelled)
		return nil
	}
	if errors.Is(err, environment.ErrAlreadyConfigured) {
		printAlreadyConfigured(out, env, deps.displayStoragePath)
		return nil
	}
	if err != nil {
		return setupError(err)
	}
	printSuccess(out, env, deps.displayStoragePath)
	return nil
}

func runNonInteractive(storagePath string, deps dependencies, out output) error {
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
	env, err := deps.setupResolved(resolved)
	if err != nil {
		return setupError(err)
	}
	printSuccess(out, env, deps.displayStoragePath)
	return nil
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
