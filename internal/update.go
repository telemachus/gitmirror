package cli

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	updateCmd     = "gitmirror-update"
	updateVersion = "v0.1.0"
	updateUsage   = `usage: gitmirror-update [-config FILENAME]

options:
    -config/-c FILENAME Specify a configuration file (default = ~/.gitmirror.json)
    -quiet/-q           Suppress output unless an error occurs
    -help/-h            Show this message
    -version/-v         Show version`
)

// Update runs git repote update on a group of repositories.
func (app *App) Update(repos []Repo) {
	if app.InitWanted || app.NoOp() {
		return
	}
	ch := make(chan Publisher)
	for _, repo := range repos {
		go app.update(repo, ch)
	}
	for range repos {
		result := <-ch
		result.Publish(app.QuietWanted)
	}
}

func (app *App) update(repo Repo, ch chan<- Publisher) {
	args := []string{"remote", "update"}
	cmd := exec.Command("git", args...)
	cmd.Dir = filepath.Join(app.HomeDir, defaultStorage, repo.Name)
	cmdString := fmt.Sprintf(
		"git %s (in %s)",
		strings.Join(args, " "),
		strings.Replace(cmd.Dir, app.HomeDir, "~", 1),
	)
	err := cmd.Run()
	if err != nil {
		app.ExitValue = exitFailure
		ch <- Failure{msg: fmt.Sprintf("%s: %s: %s", appName, cmdString, err)}
		return
	}
	ch <- Success{msg: fmt.Sprintf("%s: %s", appName, cmdString)}
}

// CmdUpdate runs `git remote update` on repos listed in a config file.
func CmdUpdate(args []string) int {
	app := NewApp()
	configFile, configIsDefault := app.Flags(args)
	repos := app.Unmarshal(configFile, configIsDefault)
	app.Update(repos)
	return app.ExitValue
}
