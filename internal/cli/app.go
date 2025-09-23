// Package cli creates and runs a command line interface.
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/telemachus/opts"
)

type syncResults struct {
	cloned   []string
	updated  []string
	upToDate []string
	errors   []string
}

type cmdEnv struct {
	results       *syncResults
	name          string
	version       string
	confFile      string
	homeDir       string
	dataDir       string
	exitValue     int
	helpWanted    bool
	quietWanted   bool
	versionWanted bool
}

func cmdFrom(name, version string, args []string) (*cmdEnv, error) {
	cmd := &cmdEnv{
		name:      name,
		version:   version,
		exitValue: exitSuccess,
		results:   &syncResults{},
	}

	og := opts.NewGroup(cmd.name)
	og.StringZero(&cmd.confFile, "config")
	og.Bool(&cmd.helpWanted, "help")
	og.Bool(&cmd.helpWanted, "h")
	og.Bool(&cmd.quietWanted, "quiet")
	og.Bool(&cmd.versionWanted, "version")

	if err := og.Parse(args); err != nil {
		return cmd, fmt.Errorf("%s: %s\n%s", cmd.name, err, cmdUsage)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cmd, fmt.Errorf("%s: %s", cmd.name, err)
	}
	cmd.homeDir = homeDir

	if cmd.confFile == "" {
		cmd.confFile = filepath.Join(cmd.homeDir, confFile)
	}
	cmd.dataDir = filepath.Join(cmd.homeDir, dataDir)

	return cmd, nil
}

func (cmd *cmdEnv) printHelpOrVersion() {
	switch {
	// Do not use cmd.noOp since that checks helpWanted and versionWanted.
	case cmd.exitValue != exitSuccess:
		return
	case cmd.helpWanted:
		fmt.Println(cmdUsage)
	case cmd.versionWanted:
		fmt.Printf("%s: %s\n", cmd.name, cmd.version)
	}
}

func (cmd *cmdEnv) noOp() bool {
	return cmd.exitValue != exitSuccess || cmd.helpWanted || cmd.versionWanted
}
