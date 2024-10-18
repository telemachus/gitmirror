// Package cli organizes and implements a command line program.
package cli

const (
	appName    = "gitmirror"
	appVersion = "v0.1.0"
	appUsage   = `usage: gitmirror [-config FILENAME]

options:
    -config FILENAME    Specify a configuration file (default = ~/.gitmirror.json)
    -quiet/-q           Suppress output unless an error occurs
    -help               Show this message
    -version            Show version`
	defaultConfig = ".gitmirror.json"
	exitSuccess   = 0
	exitFailure   = 1
)

// Run creates an App, performs the App's tasks, and returns an exit value.
func Run(args []string) int {
	app := &App{ExitValue: exitSuccess}
	configFile, isDefault := app.ParseFlags(args)
	repos := app.Unmarshal(configFile, isDefault)
	app.Mirror(repos)
	return app.ExitValue
}
