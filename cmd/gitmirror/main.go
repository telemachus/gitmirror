// gitmirror uses git commands to backup git repositories.
package main

import (
	"os"

	"github.com/telemachus/gitmirror/internal/cli"
)

func main() {
	os.Exit(cli.Gitmirror(os.Args))
}
