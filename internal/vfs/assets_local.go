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

//go:build !embed
// +build !embed

package vfs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	localAssetDirs []string
)

// initAssets initializes the VFS. Must run after loading the assets.
func initAssets() error {
	localAssetDirs = []string{"assets"}
	content, err := ioutil.ReadDir("third_party")
	if err != nil {
		return fmt.Errorf("could not find local third party directory: %v", err)
	}
	for _, info := range content {
		if !info.IsDir() {
			continue
		}
		localAssetDirs = append(localAssetDirs, filepath.Join("third_party", info.Name(), "assets"))
	}
	log.Debugf("local asset search path: %v", localAssetDirs)
	return nil
}

// load loads a file from the VFS.
func load(vfsPath string) (ReadSeekCloser, error) {
	var err error
	for _, dir := range localAssetDirs {
		var r ReadSeekCloser
		r, err = os.Open(path.Join(dir, vfsPath))
		if err != nil {
			continue
		}
		return r, nil
	}
	return nil, fmt.Errorf("could not open local:%v: %w", vfsPath, err)
}

// readDir lists all files in a directory. Returns their VFS paths!
func readDir(vfsPath string) ([]string, error) {
	var results []string
	for _, dir := range localAssetDirs {
		content, err := ioutil.ReadDir(path.Join(dir, vfsPath))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("could not scan local:%v:%v: %v", vfsPath, dir, err)
			}
			continue
		}
		for _, info := range content {
			results = append(results, path.Join(vfsPath, info.Name()))
		}
	}
	sort.Strings(results)
	return results, nil
}
