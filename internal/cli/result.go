package cli

import (
	"fmt"
	"os"
)

type result struct {
	msg   string
	isErr bool
}

func (r result) publish(quiet bool) {
	if r.isErr {
		fmt.Fprintln(os.Stderr, r.msg)
		return
	}
	if quiet {
		return
	}
	fmt.Fprintln(os.Stdout, r.msg)
}
