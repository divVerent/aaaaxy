package vfs

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/rakyll/statik/fs"

	_ "github.com/divVerent/aaaaaa/internal/assets/statik"
)

var (
	myfs http.FileSystem

	assetsDir = flag.String("assets_dir", "", "if set, use this directory as asset source instead of built-in")
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
	vfsPath := path.Join("/", purpose, path.Base(name))
	if *assetsDir != "" {
		r, err := os.Open(path.Join(*assetsDir, vfsPath))
		if err != nil {
			return nil, fmt.Errorf("could not open local:%v%v: %v", *assetsDir, vfsPath, err)
		}
		return r, nil
	}
	r, err := myfs.Open(vfsPath)
	if err != nil {
		return nil, fmt.Errorf("could not open statik:%v: %v", vfsPath, err)
	}
	return r, nil
}
