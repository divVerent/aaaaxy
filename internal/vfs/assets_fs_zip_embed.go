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

//go:build zip && embed
// +build zip,embed

package vfs

import (
	"bytes"
	"io/fs"
	"time"

	"github.com/divVerent/aaaaxy/assets"
)

type fileInfo int64

func (i fileInfo) Name() string       { return "asset" }
func (i fileInfo) Size() int64        { return int64(i) }
func (i fileInfo) Mode() fs.FileMode  { return 0555 }
func (i fileInfo) ModTime() time.Time { return time.Time{} }
func (i fileInfo) IsDir() bool        { return false }
func (i fileInfo) Sys() interface{}   { return nil }

type reader struct {
	*bytes.Reader
}

/*
func (r *reader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.f.Close()
}
*/

func (r *reader) Stat() (fs.FileInfo, error) {
	return fileInfo(r.Size()), nil
}

func (r *reader) Close() error {
	return nil
}

func openAssetsZip() (reader, error) {
	return reader{
		Reader: bytes.NewReader(assets.Data),
	}, nil
}
