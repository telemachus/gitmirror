package main

import (
	"os"

	"github.com/telemachus/gitmirror/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
