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

//go:build !embed && !zip
// +build !embed,!zip

package vfs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// initAssets initializes the VFS.
func initAssets() error {
	dirs := []string{"assets"}
	content, err := ioutil.ReadDir("third_party")
	if err != nil {
		return fmt.Errorf("could not find local third party directory: %v", err)
	}
	for _, info := range content {
		if !info.IsDir() {
			continue
		}
		dirs = append(dirs, filepath.Join("third_party", info.Name(), "assets"))
	}
	return initLocalAssets(dirs)
}

// load loads a file from the VFS.
func load(vfsPath string) (ReadSeekCloser, error) {
	return loadLocal(vfsPath)
}

// readDir lists all files in a directory. Returns their VFS names, NOT full paths!
func readDir(vfsPath string) ([]string, error) {
	return readLocalDir(vfsPath)
}
