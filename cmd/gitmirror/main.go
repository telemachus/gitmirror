// gitmirror uses git itself to back up git repositories.
//
// Briefly stated, gitmirror uses two commands.
//
//	git clone --mirror <gitURL>
//	git remote update <gitURL>
//
// The first command creates an initial archive of the repository; the second
// brings the archive up to date.
package main

import (
	"os"

	"github.com/telemachus/gitmirror/internal/cli"
)

func main() {
	os.Exit(cli.Gitmirror(os.Args[1:]))
}
