package cli

import (
	"fmt"
	"os"
)

const (
	cmd         = "gitmirror"
	cmdVersion  = "v0.9.5"
	config      = ".gitmirror.json"
	storage     = ".local/share/gitmirror"
	exitSuccess = 0
	exitFailure = 1
)

// Gitmirror runs a subcommand and returns success or failure to the shell.
func Gitmirror(args []string) int {
	app, err := appFrom(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", cmd, err)
		return exitFailure
	}

	switch app.subCmd {
	case "update", "up":
		subCmdUpdate(app)
	case "clone":
		subCmdClone(app)
	case "sync":
		subCmdSync(app)
	default:
		fmt.Fprintf(os.Stderr, "%s: unrecognized subcommand %q\n", app.cmd, app.subCmd)
		app.exitVal = exitFailure
	}

	return app.exitVal
}
