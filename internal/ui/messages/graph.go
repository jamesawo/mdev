package messages

const (
	GraphCmdShortDescription = "Show the dependency graph of supported tools"
	GraphCmdLongDescription  = `Graph displays the dependency relationships between development
tools managed by mdev.

Some tools depend on others to function correctly. For example,
Java requires SDKMAN to install and manage versions, while tools
like Maven and Gradle require Java to be present before they can
be installed.

This command prints the dependency graph so you can understand the
order in which tools will be installed when running:

  mdev install --all

The output shows each tool and the tools it depends on.

Example:

  mdev graph

Possible output:

  java -> sdkman
  maven -> java
  gradle -> java
  nvm
  podman
`
	GraphTitle      = "Tool dependency graph"
	GraphIndentMark = "└─ "
)
