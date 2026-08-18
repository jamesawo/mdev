package messages

const (
	VersionCmdShortDescription = "Show mdev version information"
	VersionCmdLongDescription  = `Display the current version of mdev and basic
information about the project.`
)

func VersionInfo(ver string) string      { return "mdev " + ver }
func VersionAuthor(author string) string { return "Created by " + author }
