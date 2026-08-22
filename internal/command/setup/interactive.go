package setup

import (
	"errors"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/interactive"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

type locationChoice struct {
	path   string
	custom bool
}

func setupInteractive() (*environment.Environment, error) {
	if env, configured, err := environment.Existing(); err != nil {
		return nil, err
	} else if configured {
		if info, statErr := os.Stat(env.StoragePath); statErr == nil && info.IsDir() {
			return env, environment.ErrAlreadyConfigured
		} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return nil, statErr
		}
		return waitForConfiguredStorage(env)
	}

	defaultLocation, err := environment.DefaultLocation()
	if err != nil {
		return nil, err
	}
	for {
		choice, err := chooseLocation(defaultLocation)
		if err != nil {
			return nil, err
		}
		if choice.custom {
			inputPath, pathErr := environment.AbsoluteInputPath(choice.path)
			if pathErr != nil {
				return nil, pathErr
			}
			if _, statErr := os.Stat(inputPath); errors.Is(statErr, os.ErrNotExist) {
				printer.Section(messages.SetupMissingLocation)
				printer.Info(inputPath)
				missingChoice, err := selectChoice([]string{messages.SetupCreateIt, messages.SetupTryAgain, messages.SetupCancel})
				if err != nil {
					return nil, err
				}
				if missingChoice == 1 {
					continue
				}
				if missingChoice != 0 {
					return nil, environment.ErrSetupCancelled
				}
			}
		}
		resolved, symlinkSource, err := environment.ResolveStoragePath(choice.path)
		if err != nil {
			printer.Fail(err.Error())
			if retry, retryErr := retryLocation(); retryErr != nil {
				return nil, retryErr
			} else if retry {
				continue
			}
			return nil, environment.ErrSetupCancelled
		}

		if symlinkSource != "" {
			printer.Section(messages.SetupSymlinkTitle)
			printer.Info(messages.SetupSymlinkResolution(symlinkSource, resolved))
			printer.Info(messages.SetupWillUse(resolved))
			choice, err := selectChoice([]string{messages.SetupContinue, messages.SetupChooseAnother, messages.SetupCancel})
			if err != nil {
				return nil, err
			}
			if choice == 1 {
				continue
			}
			if choice != 0 {
				return nil, environment.ErrSetupCancelled
			}
		}

		if info, statErr := os.Stat(resolved); statErr == nil && info.IsDir() {
			printer.Section(messages.SetupExistingTitle)
			choice, err := selectChoice([]string{messages.SetupUseExisting, messages.SetupChooseAnother, messages.SetupCancel})
			if err != nil {
				return nil, err
			}
			if choice == 1 {
				continue
			}
			if choice != 0 {
				return nil, environment.ErrSetupCancelled
			}
		}

		return environment.New(resolved), nil
	}
}

func chooseLocation(defaultLocation string) (locationChoice, error) {
	volumes, err := environment.ListExternalVolumes()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return locationChoice{}, err
	}
	options := []string{messages.SetupDefaultStorageOption(environment.DisplayPath(defaultLocation))}
	for _, volume := range volumes {
		options = append(options, messages.SetupVolumeOption(volume.Name, volume.Path, volume.Writable))
	}
	options = append(options, messages.SetupCustomStorageOption)

	for {
		index, err := interactive.RadioSelect(messages.SetupChooseDirectory, options)
		if err != nil {
			return locationChoice{}, interruptionError(err)
		}
		switch {
		case index < 0:
			return locationChoice{}, environment.ErrSetupCancelled
		case index == 0:
			return locationChoice{path: defaultLocation}, nil
		case index == len(options)-1:
			location, err := interactive.Input(messages.SetupEnterStoragePath, defaultLocation)
			if err != nil {
				return locationChoice{}, interruptionError(err)
			}
			return locationChoice{path: location, custom: true}, nil
		default:
			volume := volumes[index-1]
			if !volume.Writable {
				printer.Fail(messages.SetupReadOnly(volume.Path))
				continue
			}
			return locationChoice{path: volume.Path}, nil
		}
	}
}

func retryLocation() (bool, error) {
	choice, err := selectChoice([]string{messages.SetupChooseAnother, messages.SetupCancel})
	if err != nil {
		return false, err
	}
	return choice == 0, nil
}

func waitForConfiguredStorage(env *environment.Environment) (*environment.Environment, error) {
	for {
		printer.Section(messages.SetupStorageUnavailable)
		printer.Info(messages.SetupExpected(env.StoragePath))
		choice, err := selectChoice([]string{messages.SetupConnectRetry, messages.SetupCancel})
		if err != nil {
			return nil, err
		}
		if choice != 0 {
			return nil, environment.ErrSetupCancelled
		}
		if info, statErr := os.Stat(env.StoragePath); statErr == nil && info.IsDir() {
			return env, environment.ErrAlreadyConfigured
		}
	}
}

func selectChoice(options []string) (int, error) {
	choice, err := interactive.RadioSelect("", options)
	if err != nil {
		return -1, interruptionError(err)
	}
	return choice, nil
}

func interruptionError(err error) error {
	if errors.Is(err, terminal.InterruptErr) {
		return environment.ErrSetupCancelled
	}
	return err
}
