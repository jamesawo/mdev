package cmd

import (
	"github.com/jamesawo/mdev/internal/command/doctor"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Validate and initialize the mdev environment",
	Long: `Doctor validates that your development environment is correctly
configured for use with mdev.

It performs several checks including:

  • System prerequisites (brew, curl, etc.)
  • Existing mdev environment configuration
  • External storage availability
  • Tool installation status

If no environment has been configured yet, doctor will guide you
through the setup process and allow you to choose an external drive
where development tool data and caches will be stored.

Typical usage:

  mdev doctor

This is usually the first command you run on a new machine before
installing any development tools.`,
	Run: func(cmd *cobra.Command, args []string) {

		report, err := doctor.Run()
		if err != nil {
			printer.Fail("doctor failed")
			return
		}

		//  System Section
		printer.Section("System")

		// Attempt to fix missing prerequisites
		doctor.FixMissingPrerequisites(report.System)

		for _, s := range report.System {

			if s.Status {
				printer.Success(s.Name)
			} else {
				printer.Fail(s.Name + " missing")
			}
		}

		//  Environment Section
		printer.Section("Environment")

		for _, e := range report.Environment {

			if e.Status {
				if e.Detail != "" {
					printer.Success(e.Name + ": " + e.Detail)
				} else {
					printer.Success(e.Name)
				}
			} else {
				printer.Fail(e.Name + " not configured")
			}
		}

		//  Tools Section
		printer.Section("Tools")

		for _, t := range report.Tools {

			if t.Installed {
				printer.Success(t.Name)
				continue
			}

			if len(t.Dependencies) > 0 {
				printer.Fail(t.Name + " (requires " + t.Dependencies[0] + ")")
			} else {
				printer.Fail(t.Name)
			}
		}

		//  Next Steps
		printer.Section("Next steps")

		printer.Indent(1, "Install individual tools:")
		for _, t := range report.Tools {
			if !t.Installed {
				printer.Indent(2, "mdev install "+t.Name)
			}
		}

		printer.Blank()
		printer.Info("Install everything:")
		printer.Indent(2, "mdev install --all")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
