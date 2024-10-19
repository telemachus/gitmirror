package cli

import (
	"fmt"
	"os"
)

// Publisher is the interface that wraps the Publish method.
type Publisher interface {
	Publish(quietWanted bool)
}

// Success represents `git push --mirror` when nothing goes wrong.
type Success struct {
	msg string
}

// Failure represents `git push --mirror` when something goes wrong.
type Failure struct {
	msg string
}

// Publish prints a Success's message to stdout unless the user passed -quiet.
func (s Success) Publish(quietWanted bool) {
	if quietWanted {
		return
	}
	fmt.Fprintln(os.Stdout, s.msg)
}

// Publish unconditionally prints a Failure's message to stderr.
func (f Failure) Publish(_ bool) {
	fmt.Fprintln(os.Stderr, f.msg)
}
