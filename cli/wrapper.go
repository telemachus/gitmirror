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

// CmdWrapper turns, e.g., `gitmirror update --q` into `gitmirror-update --q`.
func CmdWrapper(args []string) int {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, wrapperUsage)
		return exitFailure
	}
	switch args[0] {
	case "initialize", "init":
		binary, lookErr := exec.LookPath("gitmirror-init")
		if lookErr != nil {
			fmt.Fprintf(os.Stdout, "%s: v%s\n", wrapperCmd, lookErr)
			return exitFailure
		}
		cmdArgs := make([]string, 0, len(args))
		cmdArgs = append(cmdArgs, "gitmirror-init")
		cmdArgs = append(cmdArgs, args[1:]...)
		env := os.Environ()
		execErr := syscall.Exec(binary, cmdArgs, env)
		if execErr != nil {
			fmt.Fprintf(os.Stdout, "%s: v%s\n", wrapperCmd, execErr)
			return exitFailure
		}
	case "update", "up":
		cmdArgs := args[1:]
		cmd := exec.Command("gitmirror-update", cmdArgs...)
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s: v%s\n", wrapperCmd, err)
			return exitFailure
		}
	case "version":
		fmt.Fprintf(os.Stdout, "%s: v%s\n", wrapperCmd, wrapperVersion)
	default:
		fmt.Fprint(os.Stderr, wrapperUsage)
		return exitFailure
	}
	return exitSuccess
}
