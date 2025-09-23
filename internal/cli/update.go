package cli

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/telemachus/gitmirror/internal/git"
)

func (cmd *cmdEnv) update(rs []Repo) {
	if cmd.noOp() {
		return
	}

	ch := make(chan result)
	for _, r := range rs {
		go cmd.updateOne(r, ch)
	}
	for range rs {
		res := <-ch
		cmd.collectResult(res)
	}
}

func (cmd *cmdEnv) updateOne(r Repo, ch chan<- result) {
	rDir := filepath.Join(cmd.dataDir, r.Name)

	fhBefore, err := git.NewFetchHead(rDir)
	// FETCH_HEAD will not exist before the first update. This is normal,
	// not an error.
	if err != nil && !os.IsNotExist(err) {
		ch <- result{repo: r.Name, kind: resultError}

		return
	}

	args := []string{"remote", "update"}
	gitCmd := exec.Command("git", args...)
	noGitPrompt := "GIT_TERMINAL_PROMPT=0"
	env := append(os.Environ(), noGitPrompt)
	gitCmd.Env = env
	gitCmd.Dir = rDir

	err = gitCmd.Run()
	if err != nil {
		ch <- result{repo: r.Name, kind: resultError}

		return
	}

	fhAfter, err := git.NewFetchHead(rDir)
	if err != nil {
		ch <- result{repo: r.Name, kind: resultError}

		return
	}

	if fhBefore.Equals(fhAfter) {
		ch <- result{repo: r.Name, kind: resultUpToDate}

		return
	}

	ch <- result{repo: r.Name, kind: resultUpdated}
}
