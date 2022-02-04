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

//go:build !windows
// +build !windows

package namedpipe

import (
	"fmt"
	"io"
	"sync"
)

type WriterAt struct {
	wr io.Writer

	mu   sync.Mutex
	cond sync.Cond
	pos  int64
}

func (w *WriterAt) Write(data []byte) (int, error) {
	return w.WriteAt(data, w.pos)
}

func (w *WriterAt) WriteAt(data []byte, pos int64) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for pos > w.pos {
		w.cond.Wait()
	}
	if pos < w.pos {
		return 0, fmt.Errorf("unsupported overlapping write in WriterAt: got offset %v, want %v", pos, w.pos)
	}
	n, err := w.wr.Write(data)
	w.pos += int64(n)
	w.cond.Broadcast()
	return n, err
}

func NewWriterAt(wr io.Writer) *WriterAt {
	w := &WriterAt{
		wr: wr,
	}
	w.cond.L = &w.mu
	return w
}

type WriteCloserAt struct {
	*WriterAt
	io.Closer
}

func NewWriteCloserAt(wrc io.WriteCloser) *WriteCloserAt {
	return &WriteCloserAt{
		WriterAt: NewWriterAt(wrc),
		Closer:   wrc,
	}
}
