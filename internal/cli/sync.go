package cli

import (
	"os"
	"path/filepath"
)

type syncList struct {
	toClone  []Repo
	toUpdate []Repo
}

func (cmd *cmdEnv) classify(repos []Repo) syncList {
	toClone := make([]Repo, 0, len(repos))
	toUpdate := make([]Repo, 0, len(repos))

	for _, repo := range repos {
		repoPath := filepath.Join(cmd.dataDir, repo.Name)
		if _, err := os.Stat(repoPath); err == nil {
			toUpdate = append(toUpdate, repo)
		} else {
			toClone = append(toClone, repo)
		}
	}

	return syncList{
		toClone:  toClone,
		toUpdate: toUpdate,
	}
}
