package cli

import (
	"flag"
	"fmt"
	"io"
)

func (app *App) ParseFlags(args []string) (string, bool) {
	if app.Err != nil {
		return "", false
	}

	flags := flag.NewFlagSet("git-backup", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var helpWanted bool
	var versionWanted bool
	var configFile string
	var isDefault bool
	flags.BoolVar(&helpWanted, "help", false, "")
	flags.BoolVar(&versionWanted, "version", false, "")
	flags.StringVar(&configFile, "", "", "")

	app.Err = flags.Parse(args)
	switch {
	case helpWanted:
		app.Info = appUsage
	case versionWanted:
		app.Info = fmt.Sprintf("%s: %s", appName, appVersion)
	case configFile == "":
		configFile = defaultConfig
		isDefault = true
	case app.Err != nil:
		app.Err = fmt.Errorf("%w\n%s", app.Err, appUsage)
		app.ExitValue = exitFailure
	}
	return configFile, isDefault
}
