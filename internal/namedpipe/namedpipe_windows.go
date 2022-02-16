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

//go:build windows
// +build windows

package namedpipe

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

type Fifo struct {
	fifoBase

	listener net.Listener
}

func New(name string, bufCount, bufSize int, timeout time.Duration) (*Fifo, error) {
	tmpPath := fmt.Sprintf("\\\\.\\pipe\\%s-%d", name, rand.Int63())
	listener, err := winio.ListenPipe(tmpPath, &winio.PipeConfig{
		SecurityDescriptor: "",
		MessageMode:        false,
		InputBufferSize:    int32(bufSize),
		OutputBufferSize:   int32(bufSize),
	})
	if err != nil {
		return nil, err
	}
	f := &Fifo{
		listener: listener,
	}
	f.start(tmpPath, bufCount, timeout, f.accept)
	return f, nil
}

func (f *Fifo) accept() (io.WriteCloser, error) {
	// NOTE: there is a race condition here; before Accept() is called, the pipe
	// is not usable, but Accept() never returns before actually having a
	// connection. Lucky we have lots of initialization between creating the
	// pipe and actually launching FFmpeg.
	pipe, err := f.listener.Accept()
	if err != nil {
		return nil, err
	}
	err = f.listener.Close()
	if err != nil {
		_ = pipe.Close()
		return nil, err
	}
	return pipe, nil
}
