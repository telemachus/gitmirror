package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/telemachus/gitmirror/internal/git"
)

func subCmdUpdate(app *appEnv) {
	rs := app.repos()
	app.update(rs)
}

func (app *appEnv) update(rs []Repo) {
	if app.noOp() {
		return
	}

	ch := make(chan result)
	for _, r := range rs {
		go app.updateOne(r, ch)
	}
	for range rs {
		res := <-ch
		res.publish(app.quiet)
	}
}

func (app *appEnv) updateOne(r Repo, ch chan<- result) {
	rDir := filepath.Join(app.storage, r.Name)
	fhBefore, err := git.NewFetchHead(filepath.Join(rDir, "FETCH_HEAD"))
	if err != nil {
		ch <- result{
			isErr: true,
			msg:   fmt.Sprintf("%s %s: %s: %s", app.cmd, app.subCmd, r.Name, err),
		}
		return
	}

	args := []string{"remote", "update"}
	cmd := exec.Command("git", args...)
	noGitPrompt := "GIT_TERMINAL_PROMPT=0"
	env := append(os.Environ(), noGitPrompt)
	cmd.Env = env
	cmd.Dir = rDir

	err = cmd.Run()
	if err != nil {
		ch <- result{
			isErr: true,
			msg:   fmt.Sprintf("%s %s: %s: %s", app.cmd, app.subCmd, r.Name, err),
		}
		return
	}

	fhAfter, err := git.NewFetchHead(filepath.Join(rDir, "FETCH_HEAD"))
	if err != nil {
		ch <- result{
			isErr: true,
			msg:   fmt.Sprintf("%s %s: %s: %s", app.cmd, app.subCmd, r.Name, err),
		}
		return
	}

	if fhBefore.Equals(fhAfter) {
		ch <- result{
			isErr: false,
			msg:   fmt.Sprintf("%s: already up-to-date", r.Name),
		}
		return
	}

	ch <- result{
		isErr: false,
		msg:   fmt.Sprintf("%s: updated", r.Name),
	}
}
