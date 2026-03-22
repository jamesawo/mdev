package cmd

import (
	"github.com/jamesawo/mdev/internal/command/doctor"
	"github.com/jamesawo/mdev/internal/ui/messages"
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
			printer.Fail(messages.DoctorFailed)
			return
		}

		// System Section
		printer.Section(messages.System)

		for _, s := range report.System {
			if s.Status {
				printer.Success(s.Name)
			} else {
				printer.Fail(s.Name + " " + messages.Missing)
			}
		}

		// Environment Section
		printer.Section(messages.Environment)

		for _, e := range report.Environment {
			if e.Status {
				if e.Detail != "" {
					printer.Success(e.Name + ": " + e.Detail)
				} else {
					printer.Success(e.Name)
				}
			} else {
				printer.Fail(messages.DoctorNotConfigured(e.Name))

				printer.Indent(2, messages.DoctorHintRunInstall)
			}
		}

		// Tools Section
		printer.Section(messages.DoctorTools)

		for _, t := range report.Tools {
			if t.Installed {
				printer.Success(t.Name)
				continue
			}

			// only show the name and not its dependencies, its noisy
			printer.Fail(t.Name)
		}

		// Next Steps
		printer.Section(messages.DoctorNextSteps)

		printer.Indent(1, messages.DoctorInstallIndividual)
		for _, t := range report.Tools {
			if !t.Installed {
				printer.Indent(2, messages.DoctorInstallTool(t.Name))
			}
		}

		printer.Blank()
		printer.Info(messages.DoctorInstallEverything)
		printer.Indent(2, messages.RootWorkflowInstallAll)

		printer.Blank()
		//todo check if there are any identified issues before showing this
		printer.Info(messages.DoctorFixHint)
		printer.Indent(2, messages.DoctorFixCmd)
		printer.Blank()
		printer.Blank()

	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVar(&fix, "fix", false, "Attempt to fix detected issues")
}
