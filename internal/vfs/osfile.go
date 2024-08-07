// Copyright 2023 Google LLC
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
	"io"
	"path/filepath"

	"github.com/divVerent/aaaaxy/internal/log"
)

type OSRoot int

const (
	WorkDir OSRoot = iota
	ExeDir
)

type readFile interface {
	io.Reader
	io.Seeker
	io.Closer
}

type writeFile interface {
	io.Writer
	io.WriterAt
	io.Closer
}

func osResolve(root OSRoot, name string) string {
	switch root {
	case ExeDir:
		return filepath.Join(exeDir, name)
	case WorkDir:
		return name
	}
	log.Fatalf("osResolve: invalid root %v", root)
	return ""
}

func OSOpen(root OSRoot, name string) (readFile, error) {
	return osOpen(osResolve(root, name))
}

func OSCreate(root OSRoot, name string) (writeFile, error) {
	return osCreate(osResolve(root, name))
}
