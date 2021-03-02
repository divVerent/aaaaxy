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

// +build statik

package vfs

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"

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

// load loads a file from the VFS.
func load(vfsPath string) (ReadSeekCloser, error) {
	r, err := myfs.Open(vfsPath)
	if err != nil {
		return nil, fmt.Errorf("could not open statik:%v: %w", vfsPath, err)
	}
	return r, nil
}

// readDir lists all files in a directory. Returns their VFS paths!
func readDir(vfsPath string) ([]string, error) {
	var results []string
	dir, err := myfs.Open(vfsPath)
	if err != nil {
		return nil, fmt.Errorf("could not scan statik:%v: %v", vfsPath, err)
	}
	content, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("could not scan statik:%v: %v", vfsPath, err)
	}
	for _, info := range content {
		results = append(results, filepath.Join(vfsPath, info.Name()))
	}
	sort.Strings(results)
	return results, nil
}
