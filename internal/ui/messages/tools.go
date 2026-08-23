package messages

const (
	ToolsInstalledSuffix               = " (installed)"
	ToolsJavaDescription               = "Java runtime (via SDKMAN)"
	ToolsNVMDescription                = "Node version manager"
	ToolsSDKMANDescription             = "Java version manager"
	ToolsMavenDescription              = "Java build automation tool"
	ToolsGradleDescription             = "Build automation tool"
	ToolsPodmanDescription             = "Container runtime with Podman Desktop"
	ToolsInvalidSDKMANIdentifier       = "invalid sdkman identifier %q"
	ToolsSDKMANCurrentNotSymlink       = "sdkman candidate %s current selection is not a symbolic link"
	ToolsSDKMANCurrentOutsideCandidate = "sdkman candidate %s current selection resolves outside managed candidate storage: %s"
	ToolsSDKMANRestoreCurrent          = "restore sdkman candidate %s current selection: %w"
)
