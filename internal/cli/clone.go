package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func subCmdClone(app *appEnv) {
	rs := app.repos()
	app.clone(rs)
}

func (app *appEnv) clone(rs []Repo) {
	if app.noOp() {
		return
	}

	err := os.MkdirAll(app.storage, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %s: %s\n", app.cmd, app.subCmd, err)
		app.exitVal = exitFailure
		return
	}

	ch := make(chan result)
	for _, r := range rs {
		go app.cloneOne(r, ch)
	}
	for range rs {
		res := <-ch
		res.publish(app.quiet)
	}
}

func (app *appEnv) cloneOne(r Repo, ch chan<- result) {
	// Normally, it is a bad idea to check whether a directory exists
	// before trying an operation.  However, this case is an exception.
	// git clone --mirror /path/to/existing/repo.git will fail with an
	// error, but for the purpose of this app, there is no error.
	// If a directory with the repo's name exists, I simply want to skip
	// that repo.
	rPath := filepath.Join(app.storage, r.Name)
	if _, err := os.Stat(rPath); err == nil {
		ch <- result{
			isErr: false,
			msg:   fmt.Sprintf("%s: already present in %s", r.Name, app.prettyPath(app.storage)),
		}
		return
	}

	args := []string{"clone", "--mirror", r.URL, r.Name}
	cmd := exec.Command("git", args...)
	noGitPrompt := "GIT_TERMINAL_PROMPT=0"
	env := append(os.Environ(), noGitPrompt)
	cmd.Env = env
	cmd.Dir = app.storage

	err := cmd.Run()
	if err != nil {
		app.exitVal = exitFailure
		ch <- result{
			isErr: true,
			msg:   fmt.Sprintf("%s %s: %s: %s", app.cmd, app.subCmd, r.Name, err),
		}
		return
	}

	ch <- result{
		isErr: false,
		msg:   fmt.Sprintf("%s: cloned", r.Name),
	}
}
