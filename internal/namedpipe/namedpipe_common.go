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

//go:build !wasm
// +build !wasm

package namedpipe

import (
	"fmt"
	"io"
	"time"
)

type fifoBase struct {
	path    string
	timeout time.Duration
	buf     chan []byte
	done    chan error
	broken  bool
	init    func() (io.WriteCloser, error)
}

func (f *fifoBase) start(path string, bufCount int, timeout time.Duration, init func() (io.WriteCloser, error)) {
	*f = fifoBase{
		path:    path,
		timeout: timeout,
		buf:     make(chan []byte, bufCount),
		done:    make(chan error),
		broken:  false,
		init:    init,
	}
	go f.run()
}

func (f *fifoBase) Path() string {
	return f.path
}

func (f *fifoBase) Write(p []byte) (int, error) {
	if f.broken {
		return 0, fmt.Errorf("named pipe %v had previous error", f.path)
	}
	select {
	case f.buf <- p:
		return len(p), nil
	case err := <-f.done:
		f.broken = true
		if err == nil {
			return 0, fmt.Errorf("named pipe %v already closed", f.path)
		}
		return 0, err
	case <-time.After(f.timeout):
		f.broken = true
		return 0, fmt.Errorf("timed out writing to named pipe %v", f.path)
	}
}

func (f *fifoBase) Close() error {
	close(f.buf)
	if f.broken {
		return fmt.Errorf("named pipe %v had previous error", f.path)
	}
	f.broken = true
	select {
	case err := <-f.done:
		return err
	case <-time.After(f.timeout):
		return fmt.Errorf("timed out writing to named pipe %v", f.path)
	}
}

func (f *fifoBase) run() {
	err := f.runInternal()
	f.done <- err
	close(f.done)
}

func (f *fifoBase) runInternal() (err error) {
	pipe, err := f.init()
	if err != nil {
		return err
	}
	defer func() {
		errC := pipe.Close()
		if err == nil {
			err = errC
		}
	}()
	for {
		data, ok := <-f.buf
		if !ok {
			return nil
		}
		_, err = pipe.Write(data)
		if err != nil {
			return err
		}
	}
}
