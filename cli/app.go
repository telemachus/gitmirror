package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"slices"
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

// Mirror runs git push --mirror on a group of repositories and displays the
// result of the mirror operation.
func (app *App) Mirror(wanted *Wanted) {
	if app.NoOp() || len(wanted.Repos) == 0 {
		return
	}
	wanted.Repos = slices.DeleteFunc(wanted.Repos, func(repo *Repo) bool {
		return repo == nil || repo.Dir == ""
	})
	ch := make(chan Publisher)
	for _, repo := range wanted.Repos {
		go mirror(repo, ch)
	}
	for range wanted.Repos {
		result := <-ch
		result.Publish()
	}
}

func mirror(repo *Repo, ch chan Publisher) {
	args := []string{"push", "--mirror", repo.Remote}
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.Dir
	cmdString := fmt.Sprintf("`git %s` (in %s)", strings.Join(args, " "), cmd.Dir)
	err := cmd.Run()
	if err != nil {
		ch <- Failure{msg: fmt.Sprintf("%s: %s: %s", appName, cmdString, err)}
	}
	ch <- Success{msg: fmt.Sprintf("%s: %s", appName, cmdString)}
}
