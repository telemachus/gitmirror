package main

import (
	"os"

	"github.com/telemachus/gitmirror/cli"
)

func main() {
	os.Exit(cli.CmdInitialize(os.Args[1:]))
}
