package messages

const (
	ListCmdShortDescription = "List all supported development tools"
	ListCmdLongDescription  = `List all development tools supported by mdev.

This command displays the tools that mdev knows how to manage.
Each tool includes a short description and can be installed,
configured, and managed through the mdev lifecycle.

Typical usage:

  mdev list
`
	ListAvailableTools = "Available tools"
)
