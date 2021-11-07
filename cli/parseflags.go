package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func (app *App) ParseFlags(args []string) (string, bool) {
	if app.ShouldNoOp() {
		return "", false
	}

	flags := flag.NewFlagSet("git-backup", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var configFile string
	var isDefault bool
	flags.BoolVar(&app.HelpWanted, "help", false, "")
	flags.BoolVar(&app.VersionWanted, "version", false, "")
	flags.StringVar(&configFile, "config", "", "")

	err := flags.Parse(args)
	switch {
	case app.HelpWanted:
		app.HelpWanted = true
		fmt.Println(appUsage)
		return "", false
	case app.VersionWanted:
		app.VersionWanted = true
		fmt.Printf("%s: %s\n", appName, appVersion)
		return "", false
	case configFile == "":
		configFile = defaultConfig
		isDefault = true
	case err != nil:
		fmt.Fprintf(os.Stderr, "%s: %s\n%s", appName, err, appUsage)
		app.ExitValue = exitFailure
	}
	return configFile, isDefault
}
