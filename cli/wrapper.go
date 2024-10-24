package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const (
	wrapperCmd     = "gitmirror"
	wrapperVersion = "0.1.0"
	wrapperUsage   = `usage: gitmirror <command> [options]

The commands are:
	initialize|init
	update|up
	version

Use "gitmirror help <command>" for more information about a command.
`
)

func wrap(longCmd string, args []string) int {
	binary, lookErr := exec.LookPath(longCmd)
	if lookErr != nil {
		fmt.Fprintf(os.Stdout, "%s: v%s\n", wrapperCmd, lookErr)
		return exitFailure
	}
	args[0] = longCmd
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		fmt.Fprintf(os.Stdout, "%s: v%s\n", wrapperCmd, execErr)
		return exitFailure
	}
	return exitSuccess
}

// CmdWrapper turns, e.g., `gitmirror update --q` into `gitmirror-update --q`.
func CmdWrapper(args []string) int {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, wrapperUsage)
		return exitFailure
	}
	exitValue := exitSuccess
	switch args[0] {
	case "initialize", "init":
		exitValue = wrap("gitmirror-init", args)
	case "update", "up":
		exitValue = wrap("gitmirror-update", args)
	case "version":
		fmt.Fprintf(os.Stdout, "%s: v%s\n", wrapperCmd, wrapperVersion)
	default:
		fmt.Fprint(os.Stderr, wrapperUsage)
		exitValue = exitFailure
	}
	return exitValue
}
