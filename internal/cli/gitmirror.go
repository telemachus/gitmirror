package cli

import (
	"fmt"
	"os"
)

const (
	cmdName     = "gitmirror"
	cmdVersion  = "v0.10.0"
	confFile    = ".gitmirror.json"
	dataDir     = ".local/share/gitmirror"
	exitSuccess = 0
	exitFailure = 1
)

// Gitmirror clones and updates repos and returns success or failure.
func Gitmirror(args []string) int {
	cmd, err := cmdFrom(cmdName, cmdVersion, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitFailure
	}

	if cmd.helpWanted || cmd.versionWanted {
		return exitSuccess
	}

	repos, err := cmd.repos()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitFailure
	}

	tasks := cmd.classify(repos)

	if err := cmd.processRepositories(tasks); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitFailure
	}

	if len(cmd.results.failed) > 0 {
		return exitFailure
	}

	return exitSuccess
}

func (cmd *cmdEnv) processRepositories(tasks syncList) error {
	reporter := newConsoleReporter(cmd.quietWanted)

	reporter.start(fmt.Sprintf("%s: processing %d repositories...", cmd.name, len(tasks.toUpdate)+len(tasks.toClone)))
	defer reporter.finish(cmd.results)

	cmd.update(tasks.toUpdate)
	if err := cmd.clone(tasks.toClone); err != nil {
		return err
	}

	return nil
}
