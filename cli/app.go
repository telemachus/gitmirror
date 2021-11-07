package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type App struct {
	ExitValue     int
	HelpWanted    bool
	VersionWanted bool
}

type Repo struct {
	Dir    string
	Remote string
}

type Wanted struct {
	Repos []*Repo
}

type repoResult struct {
	err error
	res string
}

func (app *App) ShouldNoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}

func (app *App) Unmarshal(configFile string, isDefault bool) *Wanted {
	if app.ShouldNoOp() {
		return nil
	}

	if isDefault {
		usr, err := user.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
			app.ExitValue = exitFailure
			return nil
		}
		configFile = fmt.Sprintf("%s%s%s", usr.HomeDir, string(os.PathSeparator), configFile)
	}

	blob, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
		app.ExitValue = exitFailure
		return nil
	}

	var wanted Wanted
	err = toml.Unmarshal(blob, &wanted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
		app.ExitValue = exitFailure
		return nil
	}
	return &wanted
}

func (app *App) MirrorRepos(wanted *Wanted) {
	if app.ShouldNoOp() || len(wanted.Repos) == 0 {
		return
	}

	outcomes := make(chan repoResult)
	for _, repo := range wanted.Repos {
		if repo == nil || repo.Dir == "" {
			continue
		}
		go func(r *Repo) {
			args := []string{"push", "--mirror", r.Remote}
			cmd := exec.Command("git", args...)
			cmd.Dir = r.Dir
			cmdString := fmt.Sprintf("`git %s` (in %s)", strings.Join(args, " "), cmd.Dir)
			err := cmd.Run()
			if err != nil {
				outcomes <- repoResult{err: fmt.Errorf("%s: problem with %s: %s", appName, cmdString, err)}
				return
			}
			outcomes <- repoResult{res: fmt.Sprintf("%s: %s", appName, cmdString)}
		}(repo)
	}

	for range wanted.Repos {
		r := <-outcomes
		if r.err != nil {
			fmt.Fprintln(os.Stderr, r.err)
			app.ExitValue = exitFailure
		} else {
			fmt.Println(r.res)
		}
	}
}
