package git_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/telemachus/gitmirror/internal/git"
)

// testFileReader implements git.FileReader for testing
type testFileReader struct {
	files map[string][]byte
}

func (t testFileReader) ReadFile(name string) ([]byte, error) {
	if data, exists := t.files[name]; exists {
		return data, nil
	}

	return nil, fmt.Errorf("file not found: %s", name)
}

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

			testFS := testFileReader{
				files: map[string][]byte{
					filepath.Join("repo1", "FETCH_HEAD"): tc.beforeContent,
					filepath.Join("repo2", "FETCH_HEAD"): tc.afterContent,
				},
			}

			fhBefore, err := git.NewFetchHeadWithReader("repo1", testFS)
			if err != nil {
				t.Fatalf("git.NewFetchHeadWithReader(repo1) failed: %v", err)
			}

			fhAfter, err := git.NewFetchHeadWithReader("repo2", testFS)
			if err != nil {
				t.Fatalf("git.NewFetchHeadWithReader(repo2) failed: %v", err)
			}

			got := fhBefore.Equals(fhAfter)
			if got != tc.expected {
				t.Errorf("fhBefore.Equals(fhAfter) = %v; want %v", got, tc.expected)
			}
		})
	}
}
