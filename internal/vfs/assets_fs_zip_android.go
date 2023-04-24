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

//go:build zip && android
// +build zip,android

package vfs

import (
	"io"
	"io/fs"
	"sync"
	"time"

	"golang.org/x/mobile/asset"
)

type fileInfo int64

func (i fileInfo) Name() string       { return "asset" }
func (i fileInfo) Size() int64        { return int64(i) }
func (i fileInfo) Mode() fs.FileMode  { return 0555 }
func (i fileInfo) ModTime() time.Time { return time.Time{} }
func (i fileInfo) IsDir() bool        { return false }
func (i fileInfo) Sys() interface{}   { return nil }

type reader struct {
	f   asset.File
	mu  sync.Mutex
	pos int64
}

func (r *reader) readAt(p []byte, off int64) (int, error) {
	_, err := r.f.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return r.f.Read(p)
}

func (r *reader) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, err := r.readAt(p, r.pos)
	r.pos += n
	return n, err
}

func (r *reader) ReadAt(p []byte, off int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.readAt(p, off)
}

func (r *reader) Stat() (fs.FileInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	size, err := r.f.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	return fileInfo(size), nil
}

func (r *reader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.f.Close()
}

func openAssetsZip() (*reader, error) {
	f, err := asset.Open("aaaaxy.dat")
	if err != nil {
		return nil, err
	}
	return &reader{f: f}, nil
}
