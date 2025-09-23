package cli

import (
	"fmt"
	"os"
	"os/exec"
)

func (cmd *cmdEnv) clone(rs []Repo) {
	if cmd.noOp() {
		return
	}

	err := os.MkdirAll(cmd.dataDir, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", cmd.name, err)
		cmd.exitValue = exitFailure

		return
	}

	ch := make(chan result)
	for _, r := range rs {
		go cmd.cloneOne(r, ch)
	}
	for range rs {
		res := <-ch
		cmd.collectResult(res)
	}
}

func (cmd *cmdEnv) cloneOne(r Repo, ch chan<- result) {
	args := []string{"clone", "--mirror", r.URL, r.Name}
	gitCmd := exec.Command("git", args...)
	noGitPrompt := "GIT_TERMINAL_PROMPT=0"
	env := append(os.Environ(), noGitPrompt)
	gitCmd.Env = env
	gitCmd.Dir = cmd.dataDir

	err := gitCmd.Run()
	if err != nil {
		cmd.exitValue = exitFailure
		ch <- result{repo: r.Name, kind: resultError}

		return
	}

	ch <- result{repo: r.Name, kind: resultCloned}
}
