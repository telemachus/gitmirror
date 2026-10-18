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

func cmdFrom(name, version string, args []string, stdout, stderr io.Writer) (*cmdEnv, error) {
	cmd := &cmdEnv{
		stdout:  stdout,
		stderr:  stderr,
		name:    name,
		version: version,
		results: &syncResults{},
	}

	og := opts.NewGroup(cmd.name)
	og.StringZero(&cmd.confFile, "config")
	og.Bool(&cmd.helpWanted, "help")
	og.Bool(&cmd.helpWanted, "h")
	og.Bool(&cmd.quietWanted, "quiet")
	og.Bool(&cmd.versionWanted, "version")

	if err := og.ParseAll(args); err != nil {
		return cmd, fmt.Errorf("%s: %w\n%s", cmd.name, err, cmdUsage)
	}

	if cmd.helpWanted {
		fmt.Fprintln(cmd.stdout, cmdUsage)
		return cmd, nil
	}

	if cmd.versionWanted {
		fmt.Fprintf(cmd.stdout, "%s: %s\n", cmd.name, cmd.version)
		return cmd, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cmd, fmt.Errorf("%s: %w", cmd.name, err)
	}
	cmd.homeDir = homeDir

	if cmd.confFile == "" {
		cmd.confFile = filepath.Join(cmd.homeDir, confFile)
	}
	cmd.dataDir = filepath.Join(cmd.homeDir, dataDir)

	return cmd, nil
}
