package cli

import "github.com/MakeNowJust/heredoc"

var gitmirrorUsage = heredoc.Docf(`
	Usage: gitmirror <subcommand> [options]

	Back up git repositories using git itself.

	Subcommands:
		clone			Clone repositories to mirror
		update|up		Update mirrored repositories
		help			Show detailed information about a subcommand
		version			Show version
	
	Run %[1]sgitmirror help <subcommand>%[1]s for more information about a subcommand.
`, "`")

var helpUsage = heredoc.Docf(`
	Usage: gitmirror help <subcommand>

	%[1]sgitmirror help%[1]s displays detailed information about a subcommand.

	Subcommands:
		clone			Clone repositories to mirror
		update|up		Update mirrored repositories
		help                    Show detailed information about a subcommand
		version			Show version
`, "`")

var cloneUsage = heredoc.Docf(
	`Usage: gitmirror clone [-config FILENAME]

	%[1]sgitmirror clone%[1]s runs %[1]sgit clone --mirror%[1]s on repos in a configuration file.

	Options:
		-config/-c FILENAME Specify configuration file (default ~/.gitmirror.json)
		-quiet/-q           Suppress output unless an error occurs
`, "`")

var updateUsage = heredoc.Docf(
	`Usage: gitmirror update|up [-config FILENAME]

	%[1]sgitmirror update%[1]s runs %[1]sgit remote update%[1]s on repos in a configuration file.

	Options:
		-config/-c FILENAME Specify configuration file (default ~/.gitmirror.json)
		-quiet/-q           Suppress output unless an error occurs
`, "`")

var versionUsage = heredoc.Docf(`
	Usage: gitmirror version

	%[1]sgitmirror version%[1]s displays version information for gitmirror.
`, "`")
