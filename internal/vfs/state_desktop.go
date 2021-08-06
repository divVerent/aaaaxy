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

// +build !wasm

package vfs

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ReadState loads the given state file and returns its contents.
func ReadState(kind StateKind, name string) ([]byte, error) {
	path, err := pathForRead(kind, name)
	if err != nil {
		// Remap to os.ErrNotExist so callers can deal with the error on their own.
		log.Printf("Could not find path for folder%d/%s: %v", kind, name, err)
		return nil, os.ErrNotExist
	}
	return ioutil.ReadFile(path)
}

// WriteState writes the given state file.
func WriteState(kind StateKind, name string, data []byte) error {
	path, err := pathForWrite(kind, name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}
