// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/rakyll/statik/fs"

	_ "github.com/divVerent/aaaaaa/internal/assets/statik"
)

var (
	myfs           http.FileSystem
	localAssetDirs []string
)

// Init initializes the VFS. Must run after loading the assets.
func init() {
	var err error
	myfs, err = fs.New()
	if err != nil {
		log.Panicf("Could not load statik file system: %v", err)
	}
	stamp, err := myfs.Open("/use-local-assets.stamp")
	if err == nil {
		stamp.Close()
		localAssetDirs = []string{"assets"}
		content, err := ioutil.ReadDir("third_party")
		if err != nil {
			log.Panicf("Could not find third party directory: %v", err)
		}
		for _, info := range content {
			localAssetDirs = append(localAssetDirs, filepath.Join("third_party", info.Name(), "assets"))
		}
		log.Printf("Local asset search path: %v", localAssetDirs)
	}
}

// ReadSeekCloser is a typical file interface.
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

// Load loads a file from the VFS based on the given file purpose and "name".
func Load(purpose string, name string) (ReadSeekCloser, error) {
	vfsPath := path.Join("/", purpose, path.Base(name))
	if localAssetDirs != nil {
		var err error
		for _, dir := range localAssetDirs {
			var r ReadSeekCloser
			r, err = os.Open(path.Join(dir, vfsPath))
			if err != nil {
				continue
			}
			return r, nil
		}
		return nil, fmt.Errorf("could not open local:%v: %v", vfsPath, err)
	}
	r, err := myfs.Open(vfsPath)
	if err != nil {
		return nil, fmt.Errorf("could not open statik:%v: %v", vfsPath, err)
	}
	return r, nil
}
