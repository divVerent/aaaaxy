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
	"io"
	"path"

	_ "github.com/divVerent/aaaaaa/internal/assets/statik"
)

// ReadSeekCloser is a typical file interface.
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

// Canonical returns the canonical name of the given asset.
// If Canonical returns the same string, Load will load the same asset.
func Canonical(name string) string {
	if name == "" {
		return name
	}
	return path.Base(name)
}
