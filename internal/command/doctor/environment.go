package doctor

import (
	"os"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/messages"
)

// checkEnvironment verifies the configured environment.
func checkEnvironment(reporter Reporter) []Check {

	results := []Check{}

	reporter.StartCheck(messages.EnvironmentConfiguration)
	cfg, err := config.Load()

	if err != nil {

		check := Check{
			Name:   messages.EnvironmentConfiguration,
			Status: false,
			Detail: messages.EnvironmentNotConfiguredShort,
		}
		results = append(results, check)
		reporter.EnvironmentCheck(check)

		return results
	}

	env := environment.New(cfg.ExternalDrive)

	reporter.StartCheck(messages.ExternalDrive)
	externalDrive := Check{
		Name:   messages.ExternalDrive,
		Status: true,
		Detail: env.ExternalDrive,
	}
	results = append(results, externalDrive)
	reporter.EnvironmentCheck(externalDrive)

	reporter.StartCheck(messages.DataDirectory)
	_, err = os.Stat(env.DataRoot)

	dataDirectory := Check{
		Name:   messages.DataDirectory,
		Status: err == nil,
		Detail: env.DataRoot,
	}
	results = append(results, dataDirectory)
	reporter.EnvironmentCheck(dataDirectory)

	return results
}
