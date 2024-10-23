// Package cli organizes and implements command line programs.
package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

const (
	defaultConfig = ".gitmirror.json"
	exitSuccess   = 0
	exitFailure   = 1
)

// App stores information about the application's state.
type App struct {
	HomeDir       string
	ExitValue     int
	HelpWanted    bool
	QuietWanted   bool
	VersionWanted bool
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}

// NewApp returns a new App pointer.
func NewApp() *App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
		return &App{ExitValue: exitFailure}
	}
	return &App{ExitValue: exitSuccess, HomeDir: homeDir}
}

// Repo stores information about a git repository.
type Repo struct {
	URL  string
	Name string
}

// Flags handles flags and options in my finicky way.
func (app *App) Flags(args []string) (string, bool) {
	if app.NoOp() {
		return "", false
	}
	cmdFlags := flag.NewcmdFlagset("gitmirror", flag.ContinueOnError)
	cmdFlags.SetOutput(io.Discard)
	var configFile string
	var configIsDefault bool
	cmdFlags.BoolVar(&app.HelpWanted, "help", false, "")
	cmdFlags.BoolVar(&app.HelpWanted, "h", false, "")
	cmdFlags.BoolVar(&app.QuietWanted, "quiet", false, "")
	cmdFlags.BoolVar(&app.QuietWanted, "q", false, "")
	cmdFlags.BoolVar(&app.VersionWanted, "version", false, "")
	cmdFlags.BoolVar(&app.VersionWanted, "v", false, "")
	cmdFlags.StringVar(&configFile, "config", "", "")
	cmdFlags.StringVar(&configFile, "c", "", "")
	err := cmdFlags.Parse(args)
	switch {
	// This must precede all other checks.
	case err != nil:
		fmt.Fprintf(os.Stderr, "%s: %s\n%s\n", appName, err, appUsage)
		app.ExitValue = exitFailure
	case app.HelpWanted:
		fmt.Println(appUsage)
	case app.VersionWanted:
		fmt.Printf("%s: v%s\n", appName, appVersion)
	case configFile == "":
		configFile = defaultConfig
		configIsDefault = true
	}
	return configFile, configIsDefault
}

// Unmarshal reads a configuration file and returns a slice of Repo.
func (app *App) Unmarshal(configFile string, configIsDefault bool) []Repo {
	if app.NoOp() {
		return nil
	}
	if configIsDefault {
		configFile = filepath.Join(app.HomeDir, configFile)
	}
	blob, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
		app.ExitValue = exitFailure
		return nil
	}
	repos := make([]Repo, 0, 20)
	err = json.Unmarshal(blob, &repos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
		app.ExitValue = exitFailure
		return nil
	}
	// Every repository must have a URL and a directory name.
	return slices.DeleteFunc(repos, func(repo Repo) bool {
		return repo.URL == "" || repo.Name == ""
	})
}
