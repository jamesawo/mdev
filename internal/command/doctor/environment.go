package doctor

import (
	"os"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

// checkEnvironment verifies the configured environment.
func checkEnvironment(reporter Reporter) []Check {

	results := []Check{}

	cfg, err := config.Load()

	if err != nil {

		check := Check{
			Name:   "environment",
			Status: false,
			Detail: "not configured",
		}
		results = append(results, check)
		if reporter != nil {
			reporter.EnvironmentCheck(check)
		}

		return results
	}

	env := environment.New(cfg.ExternalDrive)

	externalDrive := Check{
		Name:   "external drive",
		Status: true,
		Detail: env.ExternalDrive,
	}
	results = append(results, externalDrive)
	if reporter != nil {
		reporter.EnvironmentCheck(externalDrive)
	}

	_, err = os.Stat(env.DataRoot)

	dataDirectory := Check{
		Name:   "data directory",
		Status: err == nil,
		Detail: env.DataRoot,
	}
	results = append(results, dataDirectory)
	if reporter != nil {
		reporter.EnvironmentCheck(dataDirectory)
	}

	return results
}
