// Package cli organizes and implements a command line program.
package cli

import (
	"fmt"
	"os"
)

const (
	appName    = "gitmirror"
	appVersion = "v0.1.0"
	appUsage   = `usage: gitmirror [-config FILENAME] <command> [<args>]

option:
    -config/-c FILENAME Specify a configuration file (default = ~/.gitmirror.json)

commands:
    initialize|init     Clone repositories for later mirroring
    update              Update repositories
    help <command>      Get more information about a command
    version             Show version`
	defaultConfig  = ".gitmirror.json"
	defaultStorage = ".local/share/gitmirror"
	exitSuccess    = 0
	exitFailure    = 1
)

// Run creates an App, performs the App's tasks, and returns an exit value.
func Run(args []string) int {
	app := NewApp()
	// TODO: this is foul.
	// Also, I should return the command args that remain.
	configFile, configIsDefault := app.ParseGitmirror(args)
	if len(args) < 1 && !(app.HelpWanted || app.VersionWanted) {
		fmt.Fprintln(os.Stderr, appUsage)
		return exitFailure
	}
	repos := app.Unmarshal(configFile, configIsDefault)
	switch args[0] {
	case "help":
		// TODO: parse help for a command to show.
		// if len(args) != 3 { blow up }
		// switch args[1] { show help for initialize or update }
		fmt.Fprintln(os.Stdout, appUsage)
		app.HelpWanted = true
	case "version":
		fmt.Fprintln(os.Stdout, appVersion)
		app.VersionWanted = true
	case "initialize", "init":
		app.ParseInitialize(args[1:])
		app.Initialize(repos)
	case "update":
		app.ParseUpdate(args[1:])
		app.Update(repos)
	default:
		fmt.Fprintln(os.Stderr, appUsage)
		app.ExitValue = exitFailure
	}
	return app.ExitValue
}
