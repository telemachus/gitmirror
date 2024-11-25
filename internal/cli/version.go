package cli

import (
	"fmt"
	"os"
)

// Version displays gitmirror's version.
func (app *App) Version(args []string) {
	if app.NoOp() {
		return
	}
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "%s version: no arguments accepted\n", gitmirrorName)
		fmt.Fprint(os.Stderr, versionUsage)
		app.ExitValue = exitFailure
		return
	}
	fmt.Fprintf(os.Stdout, "%s: %s\n", gitmirrorName, gitmirrorVersion)
}

func CmdVersion(args []string) int {
	app := NewApp(gitmirrorUsage)
	app.Version(args)
	return app.ExitValue
}
