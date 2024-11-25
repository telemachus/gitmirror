package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	defaultConfig  = ".gitmirror.json"
	defaultStorage = ".local/share/gitmirror"
	exitSuccess    = 0
	exitFailure    = 1
)

// App stores information about the application's state.
type App struct {
	HomeDir     string
	CmdName     string
	Usage       string
	ExitValue   int
	HelpWanted  bool
	QuietWanted bool
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess
}

// NewApp returns a new App pointer.
func NewApp(cmdUsage string) *App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", gitmirrorName, err)
		return &App{ExitValue: exitFailure}
	}
	return &App{
		CmdName:   gitmirrorName,
		Usage:     gitmirrorUsage,
		ExitValue: exitSuccess,
		HomeDir:   homeDir,
	}
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
	cmdFlags := flag.NewFlagSet("gitmirror", flag.ContinueOnError)
	cmdFlags.SetOutput(io.Discard)
	var configFile string
	var configIsDefault bool
	cmdFlags.BoolVar(&app.HelpWanted, "help", false, "")
	cmdFlags.BoolVar(&app.HelpWanted, "h", false, "")
	cmdFlags.BoolVar(&app.QuietWanted, "quiet", false, "")
	cmdFlags.BoolVar(&app.QuietWanted, "q", false, "")
	cmdFlags.StringVar(&configFile, "config", "", "")
	cmdFlags.StringVar(&configFile, "c", "", "")
	err := cmdFlags.Parse(args)
	switch {
	// This must precede all other checks.
	case err != nil:
		fmt.Fprintf(os.Stderr, "%s: %s\n%s", app.CmdName, err, app.Usage)
		app.ExitValue = exitFailure
	case app.HelpWanted:
		fmt.Fprintf(os.Stderr, "%s: use 'help' not '-help' or '-h'\n", app.CmdName)
		fmt.Fprint(os.Stderr, app.Usage)
		app.ExitValue = exitFailure
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
		fmt.Fprintf(os.Stderr, "%s: %s\n", app.CmdName, err)
		app.ExitValue = exitFailure
		return nil
	}
	repos := make([]Repo, 0, 20)
	err = json.Unmarshal(blob, &repos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", app.CmdName, err)
		app.ExitValue = exitFailure
		return nil
	}
	// Every repository must have a URL and a directory name.
	return slices.DeleteFunc(repos, func(repo Repo) bool {
		return repo.URL == "" || repo.Name == ""
	})
}

// PrettyPath replaces a user's home directory with ~ in a string.
func (app *App) PrettyPath(s string) string {
	return strings.Replace(s, app.HomeDir, "~", 1)
}
