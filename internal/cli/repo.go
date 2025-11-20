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

func (cmd *cmdEnv) repos() ([]Repo, error) {
	conf, err := os.ReadFile(cmd.confFile)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", cmd.name, err)
	}

	cfg := struct {
		Repos []Repo `json:"repos"`
	}{
		Repos: make([]Repo, 0, 20),
	}
	err = json.Unmarshal(conf, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", cmd.name, err)
	}

	// Every repository must have a URL and a directory name.
	repos := slices.DeleteFunc(cfg.Repos, func(r Repo) bool {
		return r.URL == "" || r.Name == ""
	})

	return repos, nil
}
