package cmd

import (
	"github.com/jamesawo/mdev/internal/command/doctor"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

var fix bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Inspect system, environment, and tools",
	Long: `Analyze your system and development environment.

This command reports missing prerequisites, environment issues,
and tool installation status.

Use --fix to attempt automatic remediation.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Execute fixes if requested
		if fix {
			doctor.Fix()
			return
		}

		report, err := doctor.Run()
		if err != nil {
			printer.Fail("doctor failed")
			return
		}

		// System Section
		printer.Section("System")

		for _, s := range report.System {
			if s.Status {
				printer.Success(s.Name)
			} else {
				printer.Fail(s.Name + " missing")
			}
		}

		// Environment Section
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

		// Tools Section
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

		// Next Steps
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

		printer.Blank()
		printer.Info("Fix system issues automatically:")
		printer.Indent(2, "mdev doctor --fix")

	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVar(&fix, "fix", false, "Attempt to fix detected issues")
}
