package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/telemachus/gitmirror/internal/git"
)

// Update runs git repote update on a group of repositories.
func (app *App) Update(repos []Repo) {
	if app.NoOp() {
		return
	}
	ch := make(chan result)
	for _, repo := range repos {
		go app.update(repo, ch)
	}
	for range repos {
		res := <-ch
		res.publish(app.QuietWanted)
	}
}

func (app *App) update(repo Repo, ch chan<- result) {
	repoDir := filepath.Join(app.HomeDir, defaultStorage, repo.Name)
	fhBefore, err := git.NewFetchHead(filepath.Join(repoDir, "FETCH_HEAD"))
	if err != nil {
		ch <- result{
			isErr: true,
			msg:   fmt.Sprintf("%s: %s", repo.Name, err),
		}
		return
	}
	args := []string{"remote", "update"}
	cmd := exec.Command("git", args...)
	noGitPrompt := "GIT_TERMINAL_PROMPT=0"
	env := append(os.Environ(), noGitPrompt)
	cmd.Env = env
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		ch <- result{
			isErr: true,
			msg:   fmt.Sprintf("%s: %s", repo.Name, err),
		}
		return
	}
	fhAfter, err := git.NewFetchHead(filepath.Join(repoDir, "FETCH_HEAD"))
	if err != nil {
		ch <- result{
			isErr: true,
			msg:   fmt.Sprintf("%s: %s", repo.Name, err),
		}
		return
	}
	if fhBefore.Equals(fhAfter) {
		ch <- result{
			isErr: false,
			msg:   fmt.Sprintf("%s: already up-to-date", repo.Name),
		}
		return
	}
	ch <- result{
		isErr: false,
		msg:   fmt.Sprintf("%s: updated", repo.Name),
	}
}

// CmdUpdate runs `git remote update` on repos listed in a config file.
func CmdUpdate(args []string) int {
	app := NewApp(updateUsage)
	configFile, configIsDefault := app.Flags(args)
	repos := app.Unmarshal(configFile, configIsDefault)
	app.Update(repos)
	return app.ExitValue
}
