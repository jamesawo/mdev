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
	DoctorSystem                   = "System"
	DoctorEnvironment              = "Environment"
	DoctorEnvironmentConfiguration = "Environment configuration"
	DoctorStorageLocation          = "storage path"
	DoctorStorageDirectory         = "storage directory"
	DoctorTools                    = "Tools"
	DoctorNextSteps                = "Next steps"
	DoctorFixingPrerequisites      = "Fixing system prerequisites"
	DoctorFixingSystem             = "Fixing system"
	DoctorFixSummary               = "Summary"
)

const (
	DoctorFailed                        = "doctor failed"
	DoctorEnvironmentNotConfiguredShort = "Not configured"
	DoctorEnvironmentNotSetUp           = "Development environment is not set up"
	DoctorInstallMissingPrerequisites   = "Install missing prerequisites?"
	DoctorMissing                       = "missing"
	DoctorNothingToFix                  = "Nothing to fix"
	DoctorInstallIndividual             = "Install individual tools:"
	DoctorInstallEverything             = "Install everything:"
	DoctorFixHint                       = "To fix system issues automatically:"
	DoctorChecking                      = "checking %s"
	DoctorCheckingProgress              = "checking %s..."
	DoctorMissingCheck                  = "%s %s"
	DoctorReadinessState                = "%s %s"
	DoctorRemediating                   = "%s: %s"
	DoctorVerified                      = "%s verified"
	DoctorCheckDetail                   = "%s: %s"
	DoctorInstallAllCommand             = "mdev install --all"
	DoctorFixCommand                    = "mdev doctor --fix"
)

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
