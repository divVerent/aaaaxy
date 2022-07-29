// Copyright 2022 Google LLC
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

//go:build zip
// +build zip

package vfs

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"

	"github.com/divVerent/aaaaxy/internal/log"
)

// Make it seekable.
type seekingFS struct {
	fs.FS
}

type closableBytesReader struct {
	*bytes.Reader
	f fs.File
}

func (c closableBytesReader) Close() error {
	return nil
}

func (c closableBytesReader) Stat() (fs.FileInfo, error) {
	return c.f.Stat()
}

func makeSeekable(name string, f fs.File) (fs.File, error) {
	if _, ok := f.(ReadSeekCloser); ok {
		return f, nil
	}
	info, err := f.Stat()
	if err != nil {
		log.Errorf("failed to stat %v: %v", name, err)
		return f, nil
	}
	if info.IsDir() {
		return f, nil
	}
	c, closable := f.(io.Closer)
	if closable {
		defer c.Close()
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Errorf("failed to read %v: %v", name, err)
	}
	return closableBytesReader{bytes.NewReader(data), f}, nil
}

func (s seekingFS) Open(name string) (fs.File, error) {
	f, err := s.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return makeSeekable(name, f)
}

// initAssetsFS opens the zip file systems.
func initAssetsFS() ([]fsRoot, error) {
	zipf, err := openAssetsZip()
	if err != nil {
		return nil, fmt.Errorf("could not open aaaaxy.dat: %v", err)
	}
	zipi, err := zipf.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not stat aaaaxy.dat: %v", err)
	}
	zip, err := zip.NewReader(zipf, zipi.Size())
	if err != nil {
		return nil, fmt.Errorf("could not parse aaaaxy.dat: %v", err)
	}
	return []fsRoot{{
		name:     "dat:aaaaxy.dat",
		filesys:  seekingFS{zip},
		root:     ".",
		toPrefix: "/",
	}}, nil
}
