package git

import (
	"bytes"
	"os"
)

type FetchHead []byte

func NewFetchHead(f string) (FetchHead, error) {
	fh, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return FetchHead(fh), nil
}

func (fh FetchHead) Equals(fhOther FetchHead) bool {
	return bytes.Equal(fh, fhOther)
}
