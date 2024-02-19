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

//go:build wasm
// +build wasm

package vfs

import (
	"bytes"
	"encoding/base64"
	"strings"
	"sync"
	"syscall/js"
)

type osReader struct {
	*bytes.Reader
}

func (o osReader) Close() error {
	return nil
}

func osOpen(name string) (readFile, error) {
	var data string
	err := protectJS(func() {
		data = js.Global().Get(name).String()
	})
	if err != nil {
		return nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return osReader{Reader: bytes.NewReader(decoded)}, nil
}

type osWriter struct {
	mu   sync.Mutex
	buf  []byte
	name string
}

func (w *osWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func (w *osWriter) WriteAt(p []byte, off int64) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for off < int64(len(w.buf)) {
		w.buf = append(w.buf, 0)
	}
	prefix := w.buf[:off]
	var suffix []byte
	if pos := off + int64(len(p)); pos < int64(len(w.buf)) {
		suffix = w.buf[pos:]
	}
	w.buf = append(append(prefix, p...), suffix...)
	return len(p), nil
}

func (w *osWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	var buf strings.Builder
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	_, err := enc.Write(w.buf)
	if err != nil {
		return err
	}
	err = enc.Close()
	if err != nil {
		return err
	}
	return protectJS(func() {
		js.Global().Set(w.name, js.ValueOf(buf.String()))
	})
}

func osCreate(name string) (writeFile, error) {
	return &osWriter{name: name}, nil
}
