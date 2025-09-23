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

// Gitmirror runs a subcommand and returns success or failure to the shell.
func Gitmirror(args []string) int {
	cmd, err := cmdFrom(cmdName, cmdVersion, args)
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintln(os.Stderr, err)
	}

	cmd.printHelpOrVersion()

	repos := cmd.repos()
	tasks := cmd.classify(repos)

	cmd.processRepositories(tasks)

	return cmd.exitValue
}

func (cmd *cmdEnv) processRepositories(tasks syncList) {
	if cmd.noOp() {
		return
	}

	reporter := NewConsoleReporter(cmd.quietWanted)

	reporter.Start(fmt.Sprintf("%s: processing %d repositories...", cmd.name, len(tasks.toUpdate)+len(tasks.toClone)))

	cmd.update(tasks.toUpdate)
	cmd.clone(tasks.toClone)

	reporter.Finish(cmd.results)
}
