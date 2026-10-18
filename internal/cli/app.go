// Package cli creates and runs a command line interface.
package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/telemachus/opts"
)

type syncResults struct {
	cloned   []string
	updated  []string
	upToDate []string
	failed   []string
}

type cmdEnv struct {
	stdout        io.Writer
	stderr        io.Writer
	results       *syncResults
	name          string
	version       string
	confFile      string
	homeDir       string
	dataDir       string
	helpWanted    bool
	quietWanted   bool
	versionWanted bool
}

func (cmd *cmdEnv) parse(args []string) error {
	og := opts.NewGroup(cmd.name)
	og.StringZero(&cmd.confFile, "config")
	og.Bool(&cmd.helpWanted, "help")
	og.Bool(&cmd.helpWanted, "h")
	og.Bool(&cmd.quietWanted, "quiet")
	og.Bool(&cmd.versionWanted, "version")

	if err := og.ParseAll(args); err != nil {
		return fmt.Errorf("argument parsing error: %w", err)
	}

	return nil
}

// setPaths sets the home, config, and data paths.
func (cmd *cmdEnv) setPaths() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("%s: %w", cmd.name, err)
	}
	cmd.homeDir = homeDir

	if cmd.confFile == "" {
		cmd.confFile = filepath.Join(cmd.homeDir, confFile)
	}
	cmd.dataDir = filepath.Join(cmd.homeDir, dataDir)

	return nil
}
