// Package cli organizes and implements a command line program.
package cli

const (
	appName    = "gitmirror"
	appVersion = "v0.1.0"
	appUsage   = `usage: gitmirror [-config FILENAME]

options:
    -config/-c FILENAME Specify a configuration file (default = ~/.gitmirror.json)
    -init/-i            Clone repositories for later mirroring
    -quiet/-q           Suppress output unless an error occurs
    -help/-h            Show this message
    -version/-v         Show version`
	defaultConfig  = ".gitmirror.json"
	defaultStorage = ".local/share/gitmirror"
	exitSuccess    = 0
	exitFailure    = 1
)

// Run creates an App, performs the App's tasks, and returns an exit value.
func Run(args []string) int {
	app := NewApp()
	configFile, configIsDefault := app.ParseFlags(args)
	repos := app.Unmarshal(configFile, configIsDefault)
	if app.InitWanted {
		app.Initialize(repos)
	}
	app.Update(repos)
	return app.ExitValue
}
