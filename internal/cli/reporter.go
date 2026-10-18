package cli

import (
	"fmt"
	"io"
	"strings"
)

// consoleReporter outputs reports to w, usually a terminal.
type consoleReporter struct {
	w           io.Writer
	spinner     *spinner
	quietWanted bool
}

// newConsoleReporter returns a consoleReporter that writes to w.
func newConsoleReporter(w io.Writer, quietWanted bool) *consoleReporter {
	return &consoleReporter{w: w, quietWanted: quietWanted}
}

// start initiates reporting.
func (r *consoleReporter) start(banner string) {
	if !r.quietWanted {
		r.spinner = newSpinner(r.w)
		r.spinner.start(banner)
	}
}

// finish terminates reporting.
func (r *consoleReporter) finish(results *syncResults) {
	if r.spinner != nil {
		r.spinner.stop()
	}

	if !r.quietWanted {
		r.printSummary(results)
	}
}

func (r *consoleReporter) printSummary(results *syncResults) {
	if len(results.cloned) > 0 {
		fmt.Fprintf(r.w, "    cloned: %s\n", strings.Join(results.cloned, ", "))
	}

	if len(results.updated) > 0 {
		fmt.Fprintf(r.w, "    updated: %s\n", strings.Join(results.updated, ", "))
	}

	upToDateCount := len(results.upToDate)
	switch {
	case upToDateCount > 5:
		partialUpdates := strings.Join(results.upToDate[:3], ", ")
		fmt.Fprintf(r.w, "    already up-to-date: %s, ... (%d total)\n", partialUpdates, upToDateCount)
	case upToDateCount > 0:
		fmt.Fprintf(r.w, "    already up-to-date: %s\n", strings.Join(results.upToDate, ", "))
	}

	if len(results.failed) > 0 {
		fmt.Fprintf(r.w, "    failed: %s\n", strings.Join(results.failed, ", "))
	}
}
