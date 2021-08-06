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

// +build wasm

package vfs

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ReadState loads the given state file and returns its contents.
func ReadState(kind StateKind, name string) ([]byte, error) {
	path := fmt.Sprintf("%d/%d", kind, name)
	state := js.Global().Get("localStorage").Call("getItem", js.ValueOf(path))
	if state.IsNull() {
		return nil, os.ErrNotExist
	}
	if state.Type() != js.TypeString {
		log.Printf("Unexpected localStorage type: got %v, want string.", state.Type())
		return nil, fmt.Errorf("invalid type")
	}
	return state.String(), nil
}

// WriteState writes the given state file.
func WriteState(kind StateKind, name string, data []byte) error {
	path := fmt.Sprintf("%d/%d", kind, name)
	js.Global().Get("localStorage").Call("setItem", js.ValueOf(path), js.ValueOf(string(data)))
	return nil
}
