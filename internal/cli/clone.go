package cli

import (
	"fmt"
	"os"
	"os/exec"
)

func (cmd *cmdEnv) clone(rs []Repo) error {
	err := os.MkdirAll(cmd.dataDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("%s: %w", cmd.name, err)
	}

	ch := make(chan result, len(rs))
	for _, r := range rs {
		go cmd.cloneOne(r, ch)
	}
	for range rs {
		res := <-ch
		cmd.collectResult(res)
	}

	return nil
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
		ch <- result{repo: r.Name, kind: resultFailed}

		return
	}

	ch <- result{repo: r.Name, kind: resultCloned}
}
