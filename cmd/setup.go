package cmd

import (
	"errors"
	"fmt"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
	"github.com/spf13/cobra"
)

var setupStoragePath string

var setupCmd = &cobra.Command{
	Use:   "setup",
	Args:  cobra.NoArgs,
	Short: messages.SetupCmdShortDescription,
	Long:  messages.SetupCmdLongDescription,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("yes") || cmd.InheritedFlags().Changed("yes") {
			return fmt.Errorf("setup does not support --yes")
		}
		return nil
	},
	RunE: runSetup,
}

func runSetup(cmd *cobra.Command, _ []string) error {
	if setupStoragePath == "" {
		printer.Section(messages.SetupTitle)
		env, err := environment.SetupInteractive()
		if errors.Is(err, environment.ErrSetupCancelled) {
			printer.Info("setup cancelled.")
			return nil
		}
		if errors.Is(err, environment.ErrAlreadyConfigured) {
			printer.Info("setup is complete.")
			printer.Info("storage: " + env.StoragePath)
			return nil
		}
		if err != nil {
			return fmt.Errorf("%s: %w", messages.SetupFailed, err)
		}
		printSetupSuccess(env)
		return nil
	}

	if existing, configured, err := environment.Existing(); err != nil {
		return fmt.Errorf("%s: %w", messages.SetupFailed, err)
	} else if configured {
		return fmt.Errorf("mdev is already configured at %s; setup will not replace it", existing.StoragePath)
	}

	resolved, _, err := environment.ResolveStoragePath(setupStoragePath)
	if err != nil {
		return fmt.Errorf("%s: %w", messages.SetupFailed, err)
	}
	env, err := environment.SetupResolved(resolved)
	if err != nil {
		return fmt.Errorf("%s: %w", messages.SetupFailed, err)
	}
	printSetupSuccess(env)
	return nil
}

func printSetupSuccess(env *environment.Environment) {
	printer.Section(messages.SetupReady)
	printer.Info("storage: " + env.StoragePath)
	printer.Blank()
	printer.Info(messages.SetupNextStep)
	printer.Command("mdev list")
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVar(&setupStoragePath, "storage-path", "", messages.SetupStoragePathFlag)
}
