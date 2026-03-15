package doctor

import (
	"os"

	"github.com/jamesawo/mdev/internal/infrastructure/config"
	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

// checkEnvironment verifies the configured environment.
func checkEnvironment() []Check {

	results := []Check{}

	cfg, err := config.Load()

	if err != nil {

		results = append(results, Check{
			Name:   "environment",
			Status: false,
			Detail: "not configured",
		})

		return results
	}

	env := environment.New(cfg.ExternalDrive)

	results = append(results, Check{
		Name:   "external drive",
		Status: true,
		Detail: env.ExternalDrive,
	})

	_, err = os.Stat(env.DataRoot)

	results = append(results, Check{
		Name:   "data directory",
		Status: err == nil,
		Detail: env.DataRoot,
	})

	return results
}
