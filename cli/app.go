package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
)

// App stores information about the application's state.
type App struct {
	HomeDir       string
	ExitValue     int
	HelpWanted    bool
	QuietWanted   bool
	VersionWanted bool
}

// NewApp returns a new App pointer.
func NewApp() *App {
	return &App{ExitValue: exitSuccess}
}

// ParseFlags handles flags and options in my finicky way.
func (app *App) ParseFlags(args []string) (string, bool) {
	flags := flag.NewFlagSet("gitmirror", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configFile string
	var configIsDefault bool
	flags.BoolVar(&app.HelpWanted, "help", false, "")
	flags.BoolVar(&app.HelpWanted, "h", false, "")
	flags.BoolVar(&app.QuietWanted, "quiet", false, "")
	flags.BoolVar(&app.QuietWanted, "q", false, "")
	flags.BoolVar(&app.VersionWanted, "version", false, "")
	flags.BoolVar(&app.VersionWanted, "v", false, "")
	flags.StringVar(&configFile, "config", "", "")
	flags.StringVar(&configFile, "c", "", "")
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
		configIsDefault = true
	}
	return configFile, configIsDefault
}

// Repo stores information about a git repository.
type Repo struct {
	Dir    string
	Remote string
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}

// Unmarshal reads a configuration file and returns a slice of Repo.
func (app *App) Unmarshal(configFile string, configIsDefault bool) []Repo {
	if app.NoOp() {
		return nil
	}
	if configIsDefault {
		usr, err := user.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
			app.ExitValue = exitFailure
			return nil
		}
		app.HomeDir = usr.HomeDir
		configFile = filepath.Join(app.HomeDir, configFile)
	}
	blob, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
		app.ExitValue = exitFailure
		return nil
	}
	repos := make([]Repo, 0, 20)
	err = json.Unmarshal(blob, &repos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
		app.ExitValue = exitFailure
		return nil
	}
	// We cannot mirror a repository without a local directory and a remote.
	return slices.DeleteFunc(repos, func(repo Repo) bool {
		return repo.Dir == "" || repo.Remote == ""
	})
}

// Mirror runs git push --mirror on a group of repositories and displays the
// result of the mirror operation.
func (app *App) Mirror(repos []Repo) {
	if app.NoOp() || len(repos) == 0 {
		return
	}
	ch := make(chan Publisher)
	for _, repo := range repos {
		go app.mirror(repo, ch)
	}
	for range repos {
		result := <-ch
		result.Publish(app.QuietWanted)
	}
}

func (app *App) mirror(repo Repo, ch chan<- Publisher) {
	args := []string{"push", "--mirror", repo.Remote}
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.Dir
	cmdString := fmt.Sprintf(
		"git %s (in %s)",
		strings.Join(args, " "),
		strings.Replace(cmd.Dir, app.HomeDir, "~", 1),
	)
	err := cmd.Run()
	if err != nil {
		app.ExitValue = exitFailure
		ch <- Failure{msg: fmt.Sprintf("%s: %s: %s", appName, cmdString, err)}
		return
	}
	ch <- Success{msg: fmt.Sprintf("%s: %s", appName, cmdString)}
}
