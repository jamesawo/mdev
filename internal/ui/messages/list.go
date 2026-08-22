package messages

const (
	ListCmdShortDescription = "list supported tools and their installation status"
	ListCmdLongDescription  = `list every tool supported by mdev and its installation status.

system tools and other tools are shown separately and sorted alphabetically.
the command only observes current state; it does not install, configure,
repair, or otherwise modify tools.

typical usage:

  mdev list
  mdev list --json
`
	ListJSONFlag            = "output one machine-readable JSON document"
	ListSystemTools         = "system tools"
	ListTools               = "tools"
	ListNotConfigured       = "mdev is not configured; run mdev setup"
	ListStorageUnavailable  = "configured storage is unavailable at %s: %v"
	ListStorageNotDirectory = "configured storage is unavailable at %s: expected a directory"
	ListUnknownDetail       = "could not determine %s status: %v"
	ListUnknownStatuses     = "one or more tool statuses are unknown"
)
