package git

import "os"

// FileReader wraps the ReadFile interface.
type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type osFileReader struct{}

func (osFileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Provide os.ReadFile as a default for production use.
var defaultFileReader FileReader = osFileReader{}
