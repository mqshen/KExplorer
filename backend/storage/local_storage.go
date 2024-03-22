package storage

import (
	"github.com/vrischmann/userdir"
	"os"
	"path"
)

type localStorage struct {
	ConfPath string
}

// NewLocalStore returns a localStore instance.
func NewLocalStore(filename string) *localStorage {
	return &localStorage{
		ConfPath: path.Join(userdir.GetConfigHome(), "KExplorer", filename),
	}
}

// Load reads the given file in the user's configuration directory and returns
// its contents.
func (l *localStorage) Load() ([]byte, error) {
	d, err := os.ReadFile(l.ConfPath)
	if err != nil {
		return nil, err
	}
	return d, err
}
