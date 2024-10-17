package cli

import (
	"fmt"
	"os"
)

// Publisher is the interface that wraps the Publish method.
type Publisher interface {
	Publish(quietWanted bool)
}

type Success struct {
	msg string
}

type Failure struct {
	msg string
}

func (s Success) Publish(quietWanted bool) {
	if quietWanted {
		return
	}
	fmt.Fprintln(os.Stdout, s.msg)
}

func (f Failure) Publish(_ bool) {
	fmt.Fprintln(os.Stderr, f.msg)
}
