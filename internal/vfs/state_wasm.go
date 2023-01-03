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

//go:build wasm
// +build wasm

package vfs

import (
	"errors"
	"fmt"
	"os"
	"syscall/js"

	"github.com/divVerent/aaaaxy/internal/log"
)

func initState() error {
	log.Infof("configs will be written to localStorage['%d/*']", Config)
	log.Infof("save games will be written to localStorage['%d/*']", SavedGames)
	return nil
}

func protectJS(f func()) (err error) {
	ok := false
	defer func() {
		if !ok {
			err = fmt.Errorf("caught JS exception: %v", recover())
		}
	}()
	f()
	ok = true
	return
}

// ReadState loads the given state file and returns its contents.
func ReadState(kind StateKind, name string) ([]byte, error) {
	path := fmt.Sprintf("%d/%s", kind, name)
	var state js.Value
	err := protectJS(func() {
		state = js.Global().Get("localStorage").Call("getItem", js.ValueOf(path))
	})
	if err != nil {
		return nil, err
	}
	if state.IsNull() {
		return nil, os.ErrNotExist
	}
	if state.Type() != js.TypeString {
		log.Errorf("unexpected localStorage type: got %v, want string", state.Type())
		return nil, errors.New("invalid type")
	}
	return []byte(state.String()), nil
}

// MoveAwayState deletes a detected-to-be-broken state file so it will not be used again.
// It will also be printed to the console for debugging.
func MoveAwayState(kind StateKind, name string) error {
	data, err := ReadState(kind, name)
	path := fmt.Sprintf("%d/%s", kind, name)
	if err == nil {
		log.Errorf("deleting broken state file %s with content: %s", path, string(data))
	} else {
		log.Errorf("deleting broken state file %s with errorr: %s", path, err)
	}
	return protectJS(func() {
		js.Global().Get("localStorage").Call("removeItem", js.ValueOf(path))
	})
}

// WriteState writes the given state file.
func WriteState(kind StateKind, name string, data []byte) error {
	path := fmt.Sprintf("%d/%s", kind, name)
	return protectJS(func() {
		js.Global().Get("localStorage").Call("setItem", js.ValueOf(path), js.ValueOf(string(data)))
	})
}
