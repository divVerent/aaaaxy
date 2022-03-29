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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	localAssetDirs []string
)

// initLocalAssets initializes the local VFS.
func initLocalAssets(dirs []string) error {
	localAssetDirs = dirs
	log.Debugf("local asset search path: %v", localAssetDirs)
	return nil
}

// loadLocal loads a file from the local VFS.
func loadLocal(vfsPath string) (ReadSeekCloser, error) {
	var err error
	for _, dir := range localAssetDirs {
		var r ReadSeekCloser
		r, err = os.Open(filepath.Join(dir, filepath.FromSlash(vfsPath)))
		if err != nil {
			continue
		}
		return r, nil
	}
	return nil, fmt.Errorf("could not open local:%v: %w", vfsPath, err)
}

// readLocalDir lists all files in a directory. Returns their VFS names, NOT full paths!
func readLocalDir(vfsPath string) ([]string, error) {
	var results []string
	for _, dir := range localAssetDirs {
		content, err := ioutil.ReadDir(filepath.Join(dir, filepath.FromSlash(vfsPath)))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("could not scan local:%v:%v: %v", vfsPath, dir, err)
			}
			continue
		}
		for _, info := range content {
			results = append(results, info.Name())
		}
	}
	sort.Strings(results)
	return results, nil
}
