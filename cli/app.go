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
	Err       error
	ErrMsg    string
	ExitValue int
	Info      string
}

type Repo struct {
	Dir    string
	Remote string
}

type Wanted struct {
	Repos []*Repo
}

func (app *App) Unmarshal(configFile string, isDefault bool) *Wanted {
	if app.Err != nil {
		return nil
	}

	if isDefault {
		var usr *user.User
		usr, app.Err = user.Current()
		if app.Err != nil {
			app.Err = fmt.Errorf("%s: %w", appName, app.Err)
			app.ExitValue = exitFailure
			return nil
		}
		configFile = fmt.Sprintf("%s%s%s", usr.HomeDir, string(os.PathSeparator), configFile)
	}

	var blob []byte
	blob, app.Err = os.ReadFile(configFile)
	if app.Err != nil {
		app.Err = fmt.Errorf("%s: %w", appName, app.Err)
		app.ExitValue = exitFailure
		return nil
	}

	var wanted Wanted
	app.Err = toml.Unmarshal(blob, &wanted)
	if app.Err != nil {
		app.ExitValue = exitFailure
	}
	return &wanted
}

func (app *App) MirrorRepos(wanted *Wanted) {
	if app.Err != nil || wanted.Repos == nil {
		return
	}

	for _, repo := range wanted.Repos {
		if repo == nil || repo.Dir == "" {
			continue
		}
		args := []string{"push", "--mirror", repo.Remote}
		cmd := exec.Command("git", args...)
		cmd.Dir = repo.Dir
		cmdString := fmt.Sprintf("`git %s` (in %s)", strings.Join(args, " "), cmd.Dir)
		app.Err = cmd.Run()
		if app.Err != nil {
			app.Err = fmt.Errorf("%s: problem with %s: %w", appName, cmdString, app.Err)
			return
		}
		fmt.Printf("%s: %s\n", appName, cmdString)
	}
}

func (app *App) DisplayInfo() {
	// Do I need—or want—the or condition here?
	if app.Info == "" || app.Err != nil {
		return
	}
	fmt.Println(app.Info)
}

func (app *App) DisplayError() {
	if app.Err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, app.Err)
}
