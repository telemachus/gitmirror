package cli

import (
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

func (cmd *cmdEnv) updateOne(r Repo, ch chan<- result) {
	rDir := filepath.Join(cmd.dataDir, r.Name)

	fhBefore, err := newFetchHead(rDir)
	// FETCH_HEAD will not exist before the first update. This is normal,
	// not an error.
	if err != nil && !os.IsNotExist(err) {
		ch <- result{repo: r.Name, kind: resultFailed}

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
		ch <- result{repo: r.Name, kind: resultFailed}

		return
	}

	fhAfter, err := newFetchHead(rDir)
	if err != nil {
		ch <- result{repo: r.Name, kind: resultFailed}

		return
	}

	if fhBefore.equals(fhAfter) {
		ch <- result{repo: r.Name, kind: resultUpToDate}

		return
	}

	ch <- result{repo: r.Name, kind: resultUpdated}
}
