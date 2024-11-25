package cli

import (
	"fmt"
	"os"
)

// Help displays the help message for a command.
func (app *App) Help(args []string) {
	if app.NoOp() {
		return
	}
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s help: no subcommand given\n", gitmirrorName)
		fmt.Fprint(os.Stderr, helpUsage)
		app.ExitValue = exitFailure
		return
	}
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "%s help: too many arguments: %+v\n", gitmirrorName, args)
		fmt.Fprint(os.Stderr, helpUsage)
		app.ExitValue = exitFailure
		return
	}
	switch args[0] {
	case "clone":
		fmt.Fprint(os.Stdout, cloneUsage)
	case "up", "update":
		fmt.Fprint(os.Stdout, updateUsage)
	case "help":
		fmt.Fprint(os.Stdout, helpUsage)
	case "version":
		fmt.Fprint(os.Stdout, versionUsage)
	default:
		fmt.Fprintf(os.Stderr, "%s help: unrecognized subcommand: \"%s\"\n", gitmirrorName, args[0])
		fmt.Fprint(os.Stderr, helpUsage)
		app.ExitValue = exitFailure
	}
}

// CmdHelp clones requested repos locally for mirroring.
func CmdHelp(args []string) int {
	app := NewApp(helpUsage)
	app.Help(args)
	return app.ExitValue
}
