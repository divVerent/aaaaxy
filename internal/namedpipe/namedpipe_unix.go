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

//go:build !wasm && !windows
// +build !wasm,!windows

package namedpipe

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type Fifo struct {
	fifoBase

	parent string
}

func New(name string, bufCount, _ int, timeout time.Duration) (*Fifo, error) {
	// NOTE: using a temporary directory as there is no other race-free way to create a temporary pipe.
	tmpDir, err := ioutil.TempDir("", name+"-*")
	if err != nil {
		return nil, err
	}
	tmpPath := filepath.Join(tmpDir, "pipe")
	err = syscall.Mkfifo(tmpPath, 0600)
	if err != nil {
		return nil, err
	}
	f := &Fifo{
		parent: tmpDir,
	}
	f.start(tmpPath, bufCount, timeout, f.accept)
	return f, nil
}

func (f *Fifo) accept() (io.WriteCloser, error) {
	pipe, err := os.OpenFile(f.path, os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(f.parent)
	if err != nil {
		_ = pipe.Close()
		return nil, err
	}
	return pipe, nil
}
