// Package cli organizes and implements a command line program.
package cli

import (
	"fmt"
	"os"
)

const (
	appName        = "gitmirror"
	appVersion     = "v0.1.0"
	gitmirrorUsage = `usage: gitmirror [-config|-c FILENAME]

options:
    -config/-c FILENAME Specify a configuration file (default = ~/.gitmirror.json)
    -init|-i            Clone repositories for later mirroring
    -quiet|-q           Suppress output unless an error occurs
    -help|-h            Show this message
    -version|-v         Show version`
	updateUsage = `usage: gitmirror update [-quiet|-q]

options:
    -quiet|-q           Suppress output unless an error occurs
    -prune|-p           Prune local branches that have been pruned on the source

Update existing mirrors using git remote update`
	initializeUsage = `usage: gitmirror initialize [-quiet|-q]

option:
    -quiet|-q           Suppress output unless an error occurs

Initialize mirrors using git clone --mirror`
	defaultConfig  = ".gitmirror.json"
	defaultStorage = ".local/share/gitmirror"
	exitSuccess    = 0
	exitFailure    = 1
)

// Run creates an App, performs the App's tasks, and returns an exit value.
func Run(args []string) int {
	app := NewApp()
	subcommandArgs := app.ParseGitmirror(args)
	// If there are no args, then I should set error and HelpWanted.
	if len(subcommandArgs) > 0 {
		repos := app.Unmarshal()
		switch subcommandArgs[0] {
		case "help":
			fmt.Fprintln(os.Stdout, gitmirrorUsage)
		case "version":
			fmt.Fprintf(os.Stdout, "%s: %s\n", appName, appVersion)
		case "update":
			app.ParseUpdate(subcommandArgs[1:])
			app.Update(repos)
		case "initialize", "init":
			app.ParseInitialize(subcommandArgs[1:])
			app.Initialize(repos)
		default:
		}
	} else if !(app.HelpWanted || app.VersionWanted) {
		fmt.Fprintln(os.Stderr, gitmirrorUsage)
		app.ExitValue = exitFailure
	}
	return app.ExitValue
}
