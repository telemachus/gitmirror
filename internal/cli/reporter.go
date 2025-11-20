package cli

import (
	"fmt"
	"strings"
)

// consoleReporter outputs reports to a terminal.
type consoleReporter struct {
	spinner     *spinner
	quietWanted bool
}

// newConsoleReporter returns a ConsoleReporter.
func newConsoleReporter(quietWanted bool) *consoleReporter {
	return &consoleReporter{quietWanted: quietWanted}
}

// start initiates reporting.
func (r *consoleReporter) start(banner string) {
	if !r.quietWanted {
		r.spinner = newSpinner()
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
		fmt.Printf("    cloned: %s\n", strings.Join(results.cloned, ", "))
	}

	if len(results.updated) > 0 {
		fmt.Printf("    updated: %s\n", strings.Join(results.updated, ", "))
	}

	upToDateCount := len(results.upToDate)
	switch {
	case upToDateCount > 5:
		partialUpdates := strings.Join(results.upToDate[:3], ", ")
		fmt.Printf("    already up-to-date: %s, ... (%d total)\n", partialUpdates, upToDateCount)
	case upToDateCount > 0:
		fmt.Printf("    already up-to-date: %s\n", strings.Join(results.upToDate, ", "))
	}

	if len(results.failed) > 0 {
		fmt.Printf("    failed: %s\n", strings.Join(results.failed, ", "))
	}
}
