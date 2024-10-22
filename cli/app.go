package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// App stores information about the application's state.
type App struct {
	HomeDir       string
	ExitValue     int
	InitWanted    bool
	HelpWanted    bool
	QuietWanted   bool
	VersionWanted bool
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}

// NewApp returns a new App pointer.
func NewApp() *App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
		return &App{ExitValue: exitFailure}
	}
	return &App{ExitValue: exitSuccess, HomeDir: homeDir}
}

// ParseFlags handles flags and options in my finicky way.
func (app *App) ParseFlags(args []string) (string, bool) {
	if app.NoOp() {
		return "", false
	}
	flags := flag.NewFlagSet("gitmirror", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configFile string
	var configIsDefault bool
	flags.BoolVar(&app.HelpWanted, "help", false, "")
	flags.BoolVar(&app.HelpWanted, "h", false, "")
	flags.BoolVar(&app.InitWanted, "init", false, "")
	flags.BoolVar(&app.InitWanted, "i", false, "")
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
	URL  string
	Name string
}

// Unmarshal reads a configuration file and returns a slice of Repo.
func (app *App) Unmarshal(configFile string, configIsDefault bool) []Repo {
	if app.NoOp() {
		return nil
	}
	if configIsDefault {
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
	// We cannot mirror a repository without a URL and a directory name.
	return slices.DeleteFunc(repos, func(repo Repo) bool {
		return repo.URL == "" || repo.Name == ""
	})
}

// Initialize runs git repote update on a group of repositories.
func (app *App) Initialize(repos []Repo) {
	if app.NoOp() {
		return
	}
	err := os.MkdirAll(filepath.Join(app.HomeDir, defaultStorage), os.ModePerm)
	if err != nil {
		app.ExitValue = exitFailure
		return
	}
	ch := make(chan Publisher)
	for _, repo := range repos {
		go app.initialize(repo, ch)
	}
	for range repos {
		result := <-ch
		result.Publish(app.QuietWanted)
	}
}

func (app *App) initialize(repo Repo, ch chan<- Publisher) {
	// Normally, it is a bad idea to check whether a directory exists
	// before trying an operation.  However, this case is an exception.
	// git clone --mirror /path/to/existing/repo.git will fail with an
	// error, but for the purpose of this app, there is no error.
	// If a directory with the same name already exists, I simply want
	// to send a Success on the channel and return.
	repoPath := filepath.Join(app.HomeDir, defaultStorage, repo.Name)
	if _, err := os.Stat(repoPath); err == nil {
		ch <- Success{msg: fmt.Sprintf("%s: %s already exists", appName, repoPath)}
		return
	}
	args := []string{"clone", "--mirror", repo.URL, repo.Name}
	cmd := exec.Command("git", args...)
	cmd.Dir = filepath.Join(app.HomeDir, defaultStorage)
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

// Update runs git repote update on a group of repositories.
func (app *App) Update(repos []Repo) {
	if app.InitWanted || app.NoOp() {
		return
	}
	ch := make(chan Publisher)
	for _, repo := range repos {
		go app.update(repo, ch)
	}
	for range repos {
		result := <-ch
		result.Publish(app.QuietWanted)
	}
}

func (app *App) update(repo Repo, ch chan<- Publisher) {
	args := []string{"remote", "update"}
	cmd := exec.Command("git", args...)
	cmd.Dir = filepath.Join(app.HomeDir, defaultStorage, repo.Name)
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
