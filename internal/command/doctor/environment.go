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

	reporter.StartCheck(messages.DoctorEnvironmentConfiguration)
	cfg, err := config.Load()

	if err != nil {

		check := Check{
			Name:   messages.DoctorEnvironmentConfiguration,
			Status: false,
			Detail: messages.DoctorEnvironmentNotConfiguredShort,
		}
		results = append(results, check)
		reporter.EnvironmentCheck(check)

		return results
	}

	env := environment.New(cfg.StoragePath)
	configuration := Check{Name: messages.DoctorEnvironmentConfiguration, Status: true}
	results = append(results, configuration)
	reporter.EnvironmentCheck(configuration)

	reporter.StartCheck(messages.DoctorStorageLocation)
	storageLocation := Check{
		Name:   messages.DoctorStorageLocation,
		Status: true,
		Detail: env.StoragePath,
	}
	results = append(results, storageLocation)
	reporter.EnvironmentCheck(storageLocation)

	reporter.StartCheck(messages.DoctorStorageDirectory)
	_, err = os.Stat(env.StoragePath)

	dataDirectory := Check{
		Name:   messages.DoctorStorageDirectory,
		Status: err == nil,
		Detail: env.StoragePath,
	}
	results = append(results, dataDirectory)
	reporter.EnvironmentCheck(dataDirectory)

	return results
}
