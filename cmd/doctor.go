package cmd

import (
	"github.com/jamesawo/mdev/internal/command/doctor"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/spf13/cobra"
)

var doctorFix bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: messages.DoctorShortDescription,
	Long:  messages.DoctorLongDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return doctor.ExecuteContext(cmd.Context(), doctorFix)
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVar(&doctorFix, "fix", false, messages.DoctorFlagFix)
}
