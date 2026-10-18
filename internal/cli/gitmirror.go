package cli

import (
	"fmt"
	"io"
)

const (
	cmdName     = "gitmirror"
	cmdVersion  = "v0.10.0"
	confFile    = ".gitmirror.json"
	dataDir     = ".local/share/gitmirror"
	exitSuccess = 0
	exitFailure = 1
)

// Gitmirror clones and updates repos, writing output to stdout and errors to
// stderr, and returns success or failure.
func Gitmirror(args []string, stdout, stderr io.Writer) int {
	cmd, err := cmdFrom(cmdName, cmdVersion, args, stdout, stderr)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitFailure
	}

	if cmd.helpWanted || cmd.versionWanted {
		return exitSuccess
	}

	repos, err := cmd.repos()
	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return exitFailure
	}

	tasks := cmd.classify(repos)

	if err := cmd.processRepositories(tasks); err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return exitFailure
	}

	if len(cmd.results.failed) > 0 {
		return exitFailure
	}

	return exitSuccess
}

func (cmd *cmdEnv) processRepositories(tasks syncList) error {
	reporter := newConsoleReporter(cmd.stdout, cmd.quietWanted)

	reporter.start(fmt.Sprintf("%s: processing %d repositories...", cmd.name, len(tasks.toUpdate)+len(tasks.toClone)))
	defer reporter.finish(cmd.results)

	cmd.update(tasks.toUpdate)

	return cmd.clone(tasks.toClone)
}
