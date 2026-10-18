package cli

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func (cmd *cmdEnv) update(rs []Repo) {
	ch := make(chan result, len(rs))
	for _, r := range rs {
		go cmd.updateOne(r, ch)
	}
	for range rs {
		res := <-ch
		cmd.collectResult(res)
	}
}

func (cmd *cmdEnv) updateOne(repo Repo, ch chan<- result) {
	repoDir := filepath.Join(cmd.dataDir, repo.Name)

	fhBefore, err := newFetchHead(repoDir)
	// FETCH_HEAD will not exist before the first update. This is normal,
	// not an error.
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		ch <- result{repo: repo.Name, kind: resultFailed}

		return
	}

	args := []string{"remote", "update"}
	gitCmd := exec.Command("git", args...)
	noGitPrompt := "GIT_TERMINAL_PROMPT=0"
	env := append(os.Environ(), noGitPrompt)
	gitCmd.Env = env
	gitCmd.Dir = repoDir

	err = gitCmd.Run()
	if err != nil {
		ch <- result{repo: repo.Name, kind: resultFailed}

		return
	}

	fhAfter, err := newFetchHead(repoDir)
	if err != nil {
		ch <- result{repo: repo.Name, kind: resultFailed}

		return
	}

	if fhBefore.equals(fhAfter) {
		ch <- result{repo: repo.Name, kind: resultUpToDate}

		return
	}

	ch <- result{repo: repo.Name, kind: resultUpdated}
}
