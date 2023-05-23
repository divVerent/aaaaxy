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

//go:build embed && !zip
// +build embed,!zip

package vfs

import (
	"fmt"
	"path"

	"github.com/divVerent/aaaaxy/assets"
	"github.com/divVerent/aaaaxy/licenses"
	"github.com/divVerent/aaaaxy/third_party"
)

// initAssetsFS opens the embedded file systems.
func initAssetsFS() ([]fsRoot, error) {
	dirs := []fsRoot{
		{
			name:     "embed:assets",
			filesys:  &assets.FS,
			root:     ".",
			toPrefix: "/",
		},
	}
	content, err := third_party.FS.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("could not find embedded third party directory: %v", err)
	}
	roots := []string{}
	for _, info := range content {
		if !info.IsDir() {
			continue
		}
		root := path.Join(info.Name(), "assets")
		dirs = append(dirs, fsRoot{
			name:     "embed:third_party/" + root,
			filesys:  &third_party.FS,
			root:     root,
			toPrefix: "/",
		})
		roots = append(roots, root)
	}
	dirs = append(dirs, fsRoot{
		name:     "embed:licenses",
		filesys:  &licenses.FS,
		root:     ".",
		toPrefix: "/licenses/",
	})
	return dirs, nil
}
