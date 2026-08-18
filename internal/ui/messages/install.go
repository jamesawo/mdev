package messages

import "fmt"

const (
	InstallCmdShortDescription = "Install a tool in your local environment."
	InstallCmdLongDescription  = `
	Install a development tool into your local environment.

This command installs a supported tool and prepares it for use with the
current mdev environment configuration. The tool will be downloaded,
installed, and configured using the paths and settings defined in your
mdev environment.

Before running this command, your environment must be initialized and
validated using 'mdev doctor'. The install process depends on the
configured directories, tool paths, and system checks performed during
that step.

If the tool is already installed, the command will detect it and skip
the installation to avoid overwriting an existing setup.

Usage:
  mdev install [tool]

Arguments:
  tool    Name of the tool to install.

Behavior:
  • Validates that the environment is configured.
  • Checks whether the requested tool is supported.
  • Detects if the tool is already installed.
  • Runs the tool-specific installation process.

Examples:
  mdev install java
  mdev install gradle
  mdev install maven

Notes:
  Each tool provides its own installation logic. The command acts as a
  dispatcher that resolves the requested tool and executes its install
  routine using the current environment configuration.
	`
)

const (
	InstallAllWithToolError = "cannot use --all with a specific tool"
	InstallPlan             = "Install plan"
	InstallStart            = "Installing tools"
	InstallContinueQuestion = "Continue installation?"
	InstallSelectTools      = "Select tools to install"
	InstallCancelled        = "Installation cancelled."
	InstallAlreadyInstalled = "already installed"
)

func InstallUnknownTool(name string) string {
	return fmt.Sprintf("unknown tool: %s", name)
}
