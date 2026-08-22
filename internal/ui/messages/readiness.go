package messages

const (
	ReadinessStateReady    = "ready"
	ReadinessStateMissing  = "missing"
	ReadinessStateOutdated = "outdated"
	ReadinessStateBroken   = "broken"

	ReadinessManualRemediation             = "manual remediation required"
	ReadinessUnknownPrerequisite           = "unknown system prerequisite %s"
	ReadinessCheckError                    = "check system prerequisite %s: %w"
	ReadinessManualRemediationError        = "system prerequisite %s is %s and requires manual remediation"
	ReadinessRemediationError              = "remediate system prerequisite %s: %w"
	ReadinessVerificationError             = "verify system prerequisite %s: %w"
	ReadinessVerificationStateError        = "verify system prerequisite %s: state is %s"
	ReadinessDependencyCycleError          = "prerequisite dependency cycle at %s"
	ReadinessUnknownDependencyError        = "%s requires unknown prerequisite %s"
	ReadinessBashOutdatedError             = "installed bash is older than version 4"
	ReadinessModernBashError               = "modern bash: %w"
	ReadinessModernBashRequired            = "Bash 4 or newer is required"
	ReadinessInstallModernBash             = "install Bash 4 or newer with Homebrew"
	ReadinessHomebrewRequired              = "Homebrew is required"
	ReadinessInstallHomebrew               = "install Homebrew"
	ReadinessInstallXcodeCLI               = "install Xcode Command Line Tools"
	ReadinessXcodeInstallationStarted      = "Xcode Command Line Tools installation has started."
	ReadinessXcodeInstallationContinuation = "complete the macOS installation, then run `mdev setup` again."
	ReadinessXcodeVerificationFailed       = "Xcode Command Line Tools installation is not complete"
	ReadinessRemediationPending            = "%s remediation requires external completion"
)
