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

// +build embed

package vfs

import (
	"embed"
	"fmt"
	"path/filepath"
	"sort"
)

//go:embed _embedroot
var myfs embed.FS

// Init initializes the VFS. Must run after loading the assets.
func Init() error {
	return nil
}

// load loads a file from the VFS.
func load(vfsPath string) (ReadSeekCloser, error) {
	f, err := myfs.Open("_embedroot" + vfsPath)
	if err != nil {
		return nil, fmt.Errorf("could not open embed:%v: %w", vfsPath, err)
	}
	rsc, ok := f.(ReadSeekCloser)
	if !ok {
		info, err := f.Stat()
		if err == nil && info.IsDir() {
			return nil, fmt.Errorf("could not open embed:%v: is a directory", vfsPath)
		}
		return nil, fmt.Errorf("could not open embed:%v: internal error (go:embed doesn't yield ReadSeekCloser)", vfsPath)
	}
	return rsc, nil
}

// readDir lists all files in a directory. Returns their VFS paths!
func readDir(vfsPath string) ([]string, error) {
	content, err := myfs.ReadDir("_embedroot" + vfsPath)
	if err != nil {
		return nil, fmt.Errorf("could not scan embed:%v: %v", vfsPath, err)
	}
	var results []string
	for _, info := range content {
		results = append(results, filepath.Join(vfsPath, info.Name()))
	}
	sort.Strings(results)
	return results, nil
}
