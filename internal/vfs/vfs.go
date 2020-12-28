package vfs

import (
	"io"
	"log"
	"net/http"
	"path"

	"github.com/rakyll/statik/fs"

	_ "github.com/divVerent/aaaaaa/internal/assets/statik"
)

var (
	myfs http.FileSystem
)

// Init initializes the VFS. Must run after loading the assets.
func init() {
	var err error
	myfs, err = fs.New()
	if err != nil {
		log.Panicf("Could not load statik file system: %v", err)
	}
}

// Load loads a file from the VFS based on the given file purpose and "name".
func Load(purpose string, name string) (io.ReadCloser, error) {
	vfsPath := path.Join(purpose, path.Base(name))
	return myfs.Open(vfsPath)
}
