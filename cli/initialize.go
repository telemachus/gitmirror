// Package cli organizes and implements a command line program.
package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	initCmd     = "gitmirror-initialize"
	initVersion = "v0.1.0"
	initUsage   = `usage: gitmirror-initialize [-config FILENAME]

options:
    -config/-c FILENAME Specify a configuration file (default = ~/.gitmirror.json)
    -quiet/-q           Suppress output unless an error occurs
    -help/-h            Show this message
    -version/-v         Show version`
)

// Initialize runs git repote update on a group of repositories.
func (app *App) Initialize(repos []Repo) {
	if app.NoOp() {
		return
	}
	err := os.MkdirAll(filepath.Join(app.HomeDir, defaultStorage), os.ModePerm)
	if err != nil {
		app.ExitValue = exitFailure
		return
	}
	ch := make(chan Publisher)
	for _, repo := range repos {
		go app.initialize(repo, ch)
	}
	for range repos {
		result := <-ch
		result.Publish(app.QuietWanted)
	}
}

func (app *App) initialize(repo Repo, ch chan<- Publisher) {
	// Normally, it is a bad idea to check whether a directory exists
	// before trying an operation.  However, this case is an exception.
	// git clone --mirror /path/to/existing/repo.git will fail with an
	// error, but for the purpose of this app, there is no error.
	// If a directory with the same name already exists, I simply want
	// to send a Success on the channel and return.
	repoPath := filepath.Join(app.HomeDir, defaultStorage, repo.Name)
	if _, err := os.Stat(repoPath); err == nil {
		ch <- Success{msg: fmt.Sprintf("%s: %s already exists", app.CmdName, repoPath)}
		return
	}
	args := []string{"clone", "--mirror", repo.URL, repo.Name}
	cmd := exec.Command("git", args...)
	cmd.Dir = filepath.Join(app.HomeDir, defaultStorage)
	cmdString := fmt.Sprintf(
		"git %s (in %s)",
		strings.Join(args, " "),
		strings.Replace(cmd.Dir, app.HomeDir, "~", 1),
	)
	err := cmd.Run()
	if err != nil {
		app.ExitValue = exitFailure
		ch <- Failure{msg: fmt.Sprintf("%s: %s: %s", app.CmdName, cmdString, err)}
		return
	}
	ch <- Success{msg: fmt.Sprintf("%s: %s", app.CmdName, cmdString)}
}

// InitRun clones requested repos locally for mirroring.
func CmdInitialize(args []string) int {
	app := NewApp(initCmd, initVersion, initUsage)
	configFile, configIsDefault := app.Flags(args)
	repos := app.Unmarshal(configFile, configIsDefault)
	app.Initialize(repos)
	return app.ExitValue
}
