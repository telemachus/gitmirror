package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Clone runs git repote update on a group of repositories.
func (app *App) Clone(repos []Repo) {
	if app.NoOp() {
		return
	}
	err := os.MkdirAll(filepath.Join(app.HomeDir, defaultStorage), os.ModePerm)
	if err != nil {
		app.ExitValue = exitFailure
		return
	}
	ch := make(chan result)
	for _, repo := range repos {
		go app.clone(repo, ch)
	}
	for range repos {
		res := <-ch
		res.publish(app.QuietWanted)
	}
}

func (app *App) clone(repo Repo, ch chan<- result) {
	// Normally, it is a bad idea to check whether a directory exists
	// before trying an operation.  However, this case is an exception.
	// git clone --mirror /path/to/existing/repo.git will fail with an
	// error, but for the purpose of this app, there is no error.
	// If a directory with the repo's name exists, I simply want to send
	// a result saying that the repo exists.
	repoPath := filepath.Join(app.HomeDir, defaultStorage, repo.Name)
	storagePath := filepath.Join(app.HomeDir, defaultStorage)
	if _, err := os.Stat(repoPath); err == nil {
		prettyPath := app.PrettyPath(storagePath)
		ch <- result{
			isErr: false,
			msg:   fmt.Sprintf("%s: already present in %s", repo.Name, prettyPath),
		}
		return
	}
	args := []string{"clone", "--mirror", repo.URL, repo.Name}
	cmd := exec.Command("git", args...)
	noGitPrompt := "GIT_TERMINAL_PROMPT=0"
	env := append(os.Environ(), noGitPrompt)
	cmd.Env = env
	cmd.Dir = filepath.Join(app.HomeDir, defaultStorage)
	err := cmd.Run()
	if err != nil {
		app.ExitValue = exitFailure
		ch <- result{
			isErr: true,
			msg:   fmt.Sprintf("%s: %s", repo.Name, err),
		}
		return
	}
	ch <- result{
		isErr: false,
		msg:   fmt.Sprintf("%s: successfully cloned", repo.Name),
	}
}

// CmdClone clones requested repos locally for mirroring.
func CmdClone(args []string) int {
	app := NewApp(cloneUsage)
	configFile, configIsDefault := app.Flags(args)
	repos := app.Unmarshal(configFile, configIsDefault)
	app.Clone(repos)
	return app.ExitValue
}
