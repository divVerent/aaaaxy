package vfs

import (
	"filepath"
	"io"

	"github.com/rakyll/statik/fs"
)

// LoadFile loads a file from the VFS based on the given file purpose and "name".
func LoadFile(purpose string, name string) (io.Reader, err) {
	vfsPath := filepath.Join(purpose, filepath.Basename(name))
	return statik.Open(vfsPath)
}
