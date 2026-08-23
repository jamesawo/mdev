package messages

import "fmt"

const (
	DoctorShortDescription = "Inspect system, environment, and tools"
	DoctorLongDescription  = `Analyze your system and development environment.

This command reports missing prerequisites, environment issues,
and tool installation status.

Use --fix to attempt automatic remediation.`
	DoctorFlagFix = "Attempt to fix detected issues"
)

const (
	DoctorSystem                   = "system requirements"
	DoctorEnvironment              = "environment"
	DoctorEnvironmentConfiguration = "configuration"
	DoctorStorageLocation          = "storage path"
	DoctorStorageDirectory         = "storage directory"
	DoctorTools                    = "tools"
	DoctorNextSteps                = "next steps"
	DoctorFixingSystem             = "fixing system requirements..."
	DoctorFixSummary               = "Summary"
)

const (
	DoctorFailed                        = "doctor failed"
	DoctorEnvironmentNotConfiguredShort = "not configured"
	DoctorEnvironmentNotSetUp           = "development environment is not set up"
	DoctorInstallMissingPrerequisites   = "allow mdev to apply these changes?"
	DoctorMissing                       = "missing"
	DoctorNothingToFix                  = "nothing to fix."
	DoctorSystemReady                   = "system requirements are ready."
	DoctorInstallIndividual             = "install individual tools:"
	DoctorInstallEverything             = "install everything:"
	DoctorFixHint                       = "to fix system issues automatically:"
	DoctorChecking                      = "checking %s"
	DoctorCheckingSystem                = "checking system requirements..."
	DoctorNeedsAttention                = "needs attention"
	DoctorSystemChangesOne              = "1 system change can be fixed by mdev."
	DoctorSystemChangesMany             = "%d system changes can be fixed by mdev."
	DoctorReadyResult                   = "✓ %s\n"
	DoctorAttentionItem                 = "  %s\n    %s\n"
	DoctorFixingItem                    = "fixing %s...\n"
	DoctorInstalledTool                 = "✓ %-*s  installed"
	DoctorMissingTool                   = "○ %-*s  missing"
	DoctorCheckingProgress              = "checking %s..."
	DoctorMissingCheck                  = "%s %s"
	DoctorReadinessState                = "%s %s"
	DoctorRemediating                   = "%s: %s"
	DoctorVerified                      = "%s verified"
	DoctorCheckDetail                   = "%s: %s"
	DoctorInstallAllCommand             = "mdev install --all"
	DoctorFixCommand                    = "mdev doctor --fix"
)

func DoctorSystemChanges(count int) string {
	if count == 1 {
		return DoctorSystemChangesOne
	}
	return fmt.Sprintf(DoctorSystemChangesMany, count)
}

func DoctorInstallationFailed(name string) string {
	return fmt.Sprintf("%s installation failed", name)
}

func DoctorInstallTool(name string) string {
	return fmt.Sprintf("mdev install %s", name)
}

func DoctorTimeElapsed(d string) string {
	return "time: " + d
}

func DoctorTotalFixTime(d string) string {
	return "Total fix time: " + d
}
