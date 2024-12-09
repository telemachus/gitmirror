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

func (app *appEnv) repos() []Repo {
	if app.noOp() {
		return nil
	}

	conf, err := os.ReadFile(app.config)
	if err != nil {
		app.exitVal = exitFailure
		fmt.Fprintf(os.Stderr, "%s %s: %s\n", app.cmd, app.subCmd, err)
		return nil
	}

	rs := make([]Repo, 0, 20)
	err = json.Unmarshal(conf, &rs)
	if err != nil {
		app.exitVal = exitFailure
		fmt.Fprintf(os.Stderr, "%s %s: %s\n", app.cmd, app.subCmd, err)
		return nil
	}

	// Every repository must have a URL and a directory name.
	return slices.DeleteFunc(rs, func(r Repo) bool {
		return r.URL == "" || r.Name == ""
	})
}
