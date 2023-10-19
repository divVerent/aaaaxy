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

package vfs

import (
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

type StateKind int

type readonlyKey struct {
	kind StateKind
	name string
}

const (
	gameName = "AAAAXY"

	Config StateKind = iota
	SavedGames
)

var (
	readonly = flag.Bool("readonly", false, "if set, save games and config changes will not be written")
)

var (
	preventWrite   *string = nil
	readonlyBuffer         = map[readonlyKey][]byte{}
)

// PreventWrite prevents further writing to any state.
//
// This is used as a safety mechanism so demo playback cannot have any
// influence on the system.
func PreventWrite(reason string) {
	preventWrite = &reason
}

// ReadState loads the given state file and returns its contents.
func ReadState(kind StateKind, name string) ([]byte, error) {
	if *readonly {
		key := readonlyKey{kind: kind, name: name}
		buf, found := readonlyBuffer[key]
		if found {
			log.Infof("readonly: forcing read of %v from memory", key)
			return append([]byte(nil), buf...), nil
		}
	}
	return readState(kind, name)
}

// WriteState writes the given state file.
func WriteState(kind StateKind, name string, data []byte) error {
	if preventWrite != nil {
		log.Fatalf("attempted to write data despite %s", *preventWrite)
	}
	if *readonly {
		key := readonlyKey{kind: kind, name: name}
		log.Infof("readonly: forcing write of %v to memory", key)
		readonlyBuffer[key] = append([]byte(nil), data...)
		return nil
	}
	return writeState(kind, name, data)
}
