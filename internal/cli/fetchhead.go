package cli

import (
	"bytes"
	"os"
	"path/filepath"
)

// fetchHead represents a git FETCH_HEAD as a []byte.
type fetchHead []byte

// newFetchHead returns the fetchHead for a git repository. It returns an error
// if the fetchHead cannot be determined.
func newFetchHead(repo string) (fetchHead, error) {
	fetchHeadPath := filepath.Join(repo, "FETCH_HEAD")
	fh, err := os.ReadFile(fetchHeadPath)
	if err != nil {
		return nil, err
	}

	return fetchHead(fh), nil
}

// equals checks whether one fetchHead is identical to another.
func (fh fetchHead) equals(fhOther fetchHead) bool {
	return bytes.Equal(fh, fhOther)
}
