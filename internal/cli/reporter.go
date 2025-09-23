package cli

import (
	"fmt"
	"strings"
)

// ProgressReporter wraps a concrete reporter's Start and Finish methods.
type ProgressReporter interface {
	Start(banner string)
	Finish(results *syncResults)
}

// ConsoleReporter outputs reports to a terminal.
type ConsoleReporter struct {
	spinner     *spinner
	quietWanted bool
}

// NewConsoleReporter returns a ConsoleReporter.
func NewConsoleReporter(quietWanted bool) *ConsoleReporter {
	return &ConsoleReporter{quietWanted: quietWanted}
}

// Start initiates reporting.
func (r *ConsoleReporter) Start(banner string) {
	if !r.quietWanted {
		r.spinner = newSpinner()
		r.spinner.start(banner)
	}
}

// Finish terminates reporting.
func (r *ConsoleReporter) Finish(results *syncResults) {
	if r.spinner != nil {
		r.spinner.stop()
	}

	if !r.quietWanted {
		r.printSummary(results)
	}
}

func (r *ConsoleReporter) printSummary(results *syncResults) {
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

	if len(results.errors) > 0 {
		fmt.Printf("    errors: %s\n", strings.Join(results.errors, ", "))
	}
}
