package main

import (
	"os"

	"github.com/telemachus/gitmirror/internal/cli"
)

func main() {
	os.Exit(cli.CmdGitmirror(os.Args[1:]))
}
