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
	"path"
	"strings"

	"github.com/divVerent/aaaaxy/internal/log"
)

// ReadSeekCloser is a typical file interface.
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

func LoadPath(purpose, name string) (ReadSeekCloser, error) {
	override := path.Base(path.Dir(name))
	if override != "" {
		purpose = override
	}
	name = path.Base(name)
	return Load(purpose, name)
}

func Load(purpose, name string) (ReadSeekCloser, error) {
	if strings.ContainsRune(name, '/') {
		log.Fatalf("noncanonical path: %v %v", purpose, name)
	}
	vfsPath := fmt.Sprintf("/%s/%s", purpose, name)
	log.Debugf("loading %v", vfsPath)
	return load(vfsPath)
}

// ReadDir lists all files in a directory. Returns their VFS paths!
func ReadDir(purpose string) ([]string, error) {
	vfsPath := fmt.Sprintf("/%s/", purpose)
	log.Debugf("listing %v", vfsPath)
	return readDir(vfsPath)
}
