package main

import (
	"os"

	"git.sr.ht/~telemachus/gitmirror/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
