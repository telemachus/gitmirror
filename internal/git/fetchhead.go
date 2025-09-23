// Package git represents and manipulates git commands and objects.
package git

import (
	"bytes"
	"path/filepath"
)

// FetchHead represents a git FETCH_HEAD as a []byte.
type FetchHead []byte

// NewFetchHead returns the FetchHead for a git repository. It returns an error
// if the FetchHead cannot be determined.
func NewFetchHead(repo string) (FetchHead, error) {
	return NewFetchHeadWithReader(repo, defaultFileReader)
}

// NewFetchHeadWithReader is like NewFetchHead but accepts a custom file reader
// for testing or other specialized situations.
func NewFetchHeadWithReader(repo string, fr FileReader) (FetchHead, error) {
	fetchHead := filepath.Join(repo, "FETCH_HEAD")
	fh, err := fr.ReadFile(fetchHead)
	if err != nil {
		return nil, err
	}

	return FetchHead(fh), nil
}

// Equals checks whether one FetchHead is identical to another.
func (fh FetchHead) Equals(fhOther FetchHead) bool {
	return bytes.Equal(fh, fhOther)
}
