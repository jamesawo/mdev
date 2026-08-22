package messages

const (
	VersionCommandName         = "version"
	VersionCmdShortDescription = "Show mdev version information"
	VersionCmdLongDescription  = `Display the identity of the running mdev binary.

This command is read-only and does not require mdev configuration or storage.`
	VersionJSONFlagName = "json"
	VersionJSONFlag     = "print version metadata as JSON"

	VersionProductLineFormat = "mdev %s\n"
	VersionCommitLineFormat  = "commit: %s\n"
	VersionBuiltLineFormat   = "built: %s\n"
	RootVersionTemplate      = "{{printf \"mdev %s\\n\" .Version}}"
	VersionWriteFailed       = "write version output: %w"
)
