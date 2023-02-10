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
	"os"
	"path/filepath"
)

// initAssets initializes the VFS.
func initAssetsFS() ([]fsRoot, error) {
	dirs := []fsRoot{
		{
			name:     "local:" + "assets",
			filesys:  os.DirFS("assets"),
			root:     ".",
			toPrefix: "/",
		},
	}
	content, err := os.ReadDir("third_party")
	if err != nil {
		return nil, fmt.Errorf("could not find local third party directory: %v", err)
	}
	for _, info := range content {
		if !info.IsDir() {
			continue
		}
		path := filepath.Join("third_party", info.Name(), "assets")
		dirs = append(dirs, fsRoot{
			name:     "local:" + path,
			filesys:  os.DirFS(path),
			root:     ".",
			toPrefix: "/",
		})
	}
	// This VFS does not support license info - which is fine as you have the
	// source and all licenses on your drive already.
	return dirs, nil
}
