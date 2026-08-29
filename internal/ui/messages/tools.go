package messages

const (
	ToolsInstalledSuffix               = " (installed)"
	ToolsJavaDescription               = "Java runtime (via SDKMAN)"
	ToolsNVMDescription                = "Node version manager"
	ToolsSDKMANDescription             = "Java version manager"
	ToolsMavenDescription              = "Java build automation tool"
	ToolsGradleDescription             = "Build automation tool"
	ToolsPodmanDescription             = "Container runtime and managed machine"
	ToolsPodmanDesktopDescription      = "Optional Podman Desktop application"
	ToolsPodmanDesktopNotInstalled     = "podman desktop application is not installed"
	ToolsPodmanUnmanagedMachine        = "podman machine %s is incomplete, but its storage ownership cannot be verified; repair or remove it manually"
	ToolsPodmanInspectOutput           = "inspect podman machine state: %w"
	ToolsPodmanReadConfig              = "read podman machine configuration: %w"
	ToolsPodmanParseConfig             = "parse podman machine configuration: %w"
	ToolsPodmanManagedStateIncomplete  = "podman machine artifacts are incomplete in managed storage"
	ToolsInvalidSDKMANIdentifier       = "invalid sdkman identifier %q"
	ToolsSDKMANCurrentNotSymlink       = "sdkman candidate %s current selection is not a symbolic link"
	ToolsSDKMANCurrentOutsideCandidate = "sdkman candidate %s current selection resolves outside managed candidate storage: %s"
	ToolsSDKMANRestoreCurrent          = "restore sdkman candidate %s current selection: %w"
)
