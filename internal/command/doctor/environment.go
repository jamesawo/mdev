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

	env := environment.New(cfg.StoragePath)

	reporter.StartCheck(messages.StorageLocation)
	storageLocation := Check{
		Name:   messages.StorageLocation,
		Status: true,
		Detail: env.StoragePath,
	}
	results = append(results, storageLocation)
	reporter.EnvironmentCheck(storageLocation)

	reporter.StartCheck(messages.StorageDirectory)
	_, err = os.Stat(env.StoragePath)

	dataDirectory := Check{
		Name:   messages.StorageDirectory,
		Status: err == nil,
		Detail: env.StoragePath,
	}
	results = append(results, dataDirectory)
	reporter.EnvironmentCheck(dataDirectory)

	return results
}
