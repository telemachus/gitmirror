package cli

import (
	"fmt"
	"os"
)

// Publisher is the interface that wraps the Publish method.
type Publisher interface {
	Publish()
}

type Success struct {
	msg string
}

type Failure struct {
	msg string
}

func (s Success) Publish() {
	fmt.Fprintln(os.Stdout, s.msg)
}

func (f Failure) Publish() {
	fmt.Fprintln(os.Stderr, f.msg)
}
