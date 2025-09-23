package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

// Repo stores information about a git repository.
type Repo struct {
	URL  string
	Name string
}

func (cmd *cmdEnv) repos() []Repo {
	if cmd.noOp() {
		return nil
	}

	conf, err := os.ReadFile(cmd.confFile)
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(os.Stderr, "%s: %s\n", cmd.name, err)

		return nil
	}

	cfg := struct {
		Repos []Repo `json:"repos"`
	}{
		Repos: make([]Repo, 0, 20),
	}
	err = json.Unmarshal(conf, &cfg)
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(os.Stderr, "%s: %s\n", cmd.name, err)

		return nil
	}

	// Every repository must have a URL and a directory name.
	return slices.DeleteFunc(cfg.Repos, func(r Repo) bool {
		return r.URL == "" || r.Name == ""
	})
}
