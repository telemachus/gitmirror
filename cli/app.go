package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// App stores information about the application's state.
type App struct {
	ExitValue     int
	HelpWanted    bool
	VersionWanted bool
}

// ParseFlags handles flags and options in my finicky way.
func (app *App) ParseFlags(args []string) (string, bool) {
	flags := flag.NewFlagSet("git-backup", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var configFile string
	var isDefault bool
	flags.BoolVar(&app.HelpWanted, "help", false, "")
	flags.BoolVar(&app.VersionWanted, "version", false, "")
	flags.StringVar(&configFile, "config", "", "")

	err := flags.Parse(args)
	switch {
	// This must precede all other checks.
	case err != nil:
		fmt.Fprintf(os.Stderr, "%s: %s\n%s\n", appName, err, appUsage)
		app.ExitValue = exitFailure
	case app.HelpWanted:
		fmt.Println(appUsage)
	case app.VersionWanted:
		fmt.Printf("%s: %s\n", appName, appVersion)
	case configFile == "":
		configFile = defaultConfig
		isDefault = true
	}
	return configFile, isDefault
}

// Repo stores information about a git repository.
type Repo struct {
	Dir    string
	Remote string
}

// Wanted stores a slice of Repo structs.
type Wanted struct {
	Repos []*Repo
}

type repoResult struct {
	err error
	res string
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}

// Unmarshal reads a configuration file and returns a Wanted struct.
func (app *App) Unmarshal(configFile string, isDefault bool) *Wanted {
	if app.NoOp() {
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

// MirrorRepos runs git push --mirror on a group of repositories.
func (app *App) MirrorRepos(wanted *Wanted) {
	if app.NoOp() || len(wanted.Repos) == 0 {
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
