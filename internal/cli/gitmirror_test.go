package cli

import (
	"bytes"
	"testing"
)

func TestHelpAndVersion(t *testing.T) {
	t.Parallel()

	testCases := map[string][]string{
		"help":    {"--help"},
		"h":       {"-h"},
		"version": {"--version"},
	}

	for msg, args := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			status := Gitmirror(args, &stdout, &stderr)
			if status != exitSuccess {
				t.Errorf("Gitmirror(%q) = %d; want %d", args, status, exitSuccess)
			}
			if stdout.Len() == 0 {
				t.Errorf("Gitmirror(%q) wrote nothing to stdout", args)
			}
		})
	}
}

func TestVersionFormat(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	Gitmirror([]string{"--version"}, &stdout, &stderr)
	if got, want := stdout.String(), cmdName+" "+cmdVersion+"\n"; got != want {
		t.Errorf("Gitmirror(--version) wrote %q; want %q", got, want)
	}
}

func TestBadOption(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	args := []string{"--no-such-option"}
	status := Gitmirror(args, &stdout, &stderr)
	if status != exitParseFailure {
		t.Errorf("Gitmirror(%q) = %d; want %d", args, status, exitParseFailure)
	}
	if stderr.Len() == 0 {
		t.Errorf("Gitmirror(%q) wrote nothing to stderr", args)
	}
}
