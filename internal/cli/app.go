// Package cli creates and runs a command line interface.
package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/telemachus/gitmirror/internal/optionparser"
)

type appEnv struct {
	cmd     string
	subCmd  string
	config  string
	home    string
	storage string
	exitVal int
	quiet   bool
}

func appFrom(args []string) (*appEnv, error) {
	app := &appEnv{cmd: cmd, exitVal: exitSuccess}

	op := optionparser.NewOptionParser()
	op.On("-c", "--config FILE", "Use FILE as config file (default ~/.gitmirror.json)", &app.config)
	op.On("-q", "--quiet", "Print only error messages", &app.quiet)
	op.On("--version", "Print version and exit", version)
	op.Command("clone", "Clone git repositories using `git clone --mirror`")
	op.Command("update|up", "Update git repositories using `git remote update`")
	op.Command("sync", "Run both update and clone (in that order)")
	op.Start = 24
	op.Banner = "Usage: gitmirror [options] <subcommand>"
	op.Coda = "\nFor more information or to file a bug report visit https://github.com/telemachus/gitmirror"

	// Do not continue if we cannot parse and validate arguments or get the
	// user's home directory.
	if err := op.ParseFrom(args); err != nil {
		return nil, err
	}
	if err := validate(op.Extra); err != nil {
		return nil, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	if app.config == "" {
		app.config = filepath.Join(home, config)
	}
	app.storage = filepath.Join(home, storage)
	app.subCmd = op.Extra[0]

	return app, nil
}

func (app *appEnv) noOp() bool {
	return app.exitVal != exitSuccess
}

func (app *appEnv) prettyPath(s string) string {
	return strings.Replace(s, app.home, "~", 1)
}

func validate(extra []string) error {
	if len(extra) != 1 {
		return errors.New("one (and only one) subcommand is required")
	}

	// The only recognized subcommands are clone, up(date), and sync.
	recognized := map[string]struct{}{
		"clone":  {},
		"update": {},
		"up":     {},
		"sync":   {},
	}
	if _, ok := recognized[extra[0]]; !ok {
		return fmt.Errorf("unrecognized subcommand: %q", extra[0])
	}

	return nil
}

// Quick and dirty, but why be fancy in this case?
func version() {
	fmt.Printf("%s %s\n", cmd, cmdVersion)
	os.Exit(exitSuccess)
}
