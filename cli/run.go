// Package cli organizes and implements a command line program.
package cli

const (
	appName    = "gitmirror"
	appVersion = "v0.0.3"
	appUsage   = `usage: gitmirror [-config FILENAME]

options:
    -config FILENAME    Specify a configuration file (default = ~/.gitmirror.toml)
    -help               Show this message
    -version            Show version`
	defaultConfig = ".gitmirror.toml"
	exitSuccess   = 0
	exitFailure   = 1
)

// Run creates an App, performs the App's tasks, and returns an exit value.
func Run(args []string) int {
	app := &App{ExitValue: exitSuccess}
	configFile, isDefault := app.ParseFlags(args)
	wanted := app.Unmarshal(configFile, isDefault)
	app.Mirror(wanted)
	return app.ExitValue
}
