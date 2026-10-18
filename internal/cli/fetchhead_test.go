package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFetchHeadEquality(t *testing.T) {
	t.Parallel()

	// Original FETCH_HEAD content
	originalContent := []byte(`089293721eb4f586907a17a18783fee1eae2f445	not-for-merge	branch 'bad-and-feel-bad' of https://github.com/owner/repo
fc558a102bc00e11580aef6033692f92d964a638	not-for-merge	branch 'clone-mirror' of https://github.com/owner/repo
aae258c89dd2c7267a88d84fe4bf1a71df274e33	not-for-merge	branch 'main' of https://github.com/owner/repo
c2647b449c1bdf91109048fe0327d738b83da1e5	not-for-merge	branch 'subcommands' of https://github.com/owner/repo
`)

	// Identical content
	identicalContent := originalContent

	// Different content (note the changed hash in the main branch line)
	differentContent := []byte(`089293721eb4f586907a17a18783fee1eae2f445	not-for-merge	branch 'bad-and-feel-bad' of https://github.com/owner/repo
fc558a102bc00e11580aef6033692f92d964a638	not-for-merge	branch 'clone-mirror' of https://github.com/owner/repo
aae458c89dd2c7267a88d84fe4bf1a71df274e33	not-for-merge	branch 'main' of https://github.com/owner/repo
c2647b449c1bdf91109048fe0327d738b83da1e5	not-for-merge	branch 'subcommands' of https://github.com/owner/repo
`)

	// Longer content (additional branch)
	longerContent := []byte(`089293721eb4f586907a17a18783fee1eae2f445	not-for-merge	branch 'bad-and-feel-bad' of https://github.com/owner/repo
fc558a102bc00e11580aef6033692f92d964a638	not-for-merge	branch 'clone-mirror' of https://github.com/owner/repo
aae258c89dd2c7267a88d84fe4bf1a71df274e33	not-for-merge	branch 'main' of https://github.com/owner/repo
497c6dbe51ac3adf1291aed2b9d6ec9de74a72e4	not-for-merge	branch 'multiple-commands' of https://github.com/owner/repo
c2647b449c1bdf91109048fe0327d738b83da1e5	not-for-merge	branch 'subcommands' of https://github.com/owner/repo
`)

	// Shorter content (missing last branch)
	shorterContent := []byte(`089293721eb4f586907a17a18783fee1eae2f445	not-for-merge	branch 'bad-and-feel-bad' of https://github.com/owner/repo
fc558a102bc00e11580aef6033692f92d964a638	not-for-merge	branch 'clone-mirror' of https://github.com/owner/repo
aae258c89dd2c7267a88d84fe4bf1a71df274e33	not-for-merge	branch 'main' of https://github.com/owner/repo
`)

	testCases := map[string]struct {
		beforeContent []byte
		afterContent  []byte
		expected      bool
	}{
		"original should equal identical": {
			beforeContent: originalContent,
			afterContent:  identicalContent,
			expected:      true,
		},
		"original should not equal different": {
			beforeContent: originalContent,
			afterContent:  differentContent,
			expected:      false,
		},
		"original should not equal longer": {
			beforeContent: originalContent,
			afterContent:  longerContent,
			expected:      false,
		},
		"original should not equal shorter": {
			beforeContent: originalContent,
			afterContent:  shorterContent,
			expected:      false,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			fhBefore := newTestFetchHead(t, tc.beforeContent)
			fhAfter := newTestFetchHead(t, tc.afterContent)

			got := fhBefore.equals(fhAfter)
			if got != tc.expected {
				t.Errorf("fhBefore.equals(fhAfter) = %v; want %v", got, tc.expected)
			}
		})
	}
}

// newTestFetchHead writes content as FETCH_HEAD in a new temporary directory
// and returns the fetchHead that newFetchHead reads from it.
func newTestFetchHead(t *testing.T, content []byte) fetchHead {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "FETCH_HEAD"), content, 0o644); err != nil {
		t.Fatalf("failed to write FETCH_HEAD: %v", err)
	}

	fh, err := newFetchHead(dir)
	if err != nil {
		t.Fatalf("newFetchHead(%s) failed: %v", dir, err)
	}

	return fh
}
