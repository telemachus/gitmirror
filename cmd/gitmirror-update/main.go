package main

import (
	"os"

	"github.com/telemachus/gitmirror/cli"
)

func main() {
	os.Exit(cli.CmdUpdate(os.Args[1:]))
}
