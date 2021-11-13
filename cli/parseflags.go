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
	// This must precede all other checks.
	case err != nil:
		fmt.Fprintf(os.Stderr, "%s: %s\n%s\n", appName, err, appUsage)
		app.ExitValue = exitFailure
	case app.HelpWanted:
		fmt.Println(appUsage)
	case app.VersionWanted:
		fmt.Printf("%s: %s\n", appName, appVersion)
	case configFile == "":
		configFile = defaultConfig
		isDefault = true
	}
	return configFile, isDefault
}
