package cli

import (
	"fmt"
	"os"
)

const (
	gitmirrorName    = "gitmirror"
	gitmirrorVersion = "v0.9.0"
)

// CmdWrapper turns, e.g., `gitmirror update -q` into `gitmirror-update -q`.
func CmdGitmirror(args []string) int {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, gitmirrorUsage)
		return exitFailure
	}
	// Give the user a hand if they try `gitmirror --help update`.
	if args[0] == "--help" || args[0] == "-help" || args[0] == "-h" {
		args[0] = "help"
	}
	var exitValue int
	switch args[0] {
	case "clone":
		exitValue = CmdClone(args[1:])
	case "update", "up":
		exitValue = CmdUpdate(args[1:])
	case "version":
		exitValue = CmdVersion(args[1:])
	case "help":
		// TODO: write the help in Asciidoc?
		exitValue = CmdHelp(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "gitmirror: unrecognized subcommand: \"%s\"\n", args[0])
		fmt.Fprint(os.Stderr, gitmirrorUsage)
		exitValue = exitFailure
	}
	return exitValue
}
