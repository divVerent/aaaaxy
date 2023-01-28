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

package vfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/divVerent/aaaaxy/internal/log"
)

func initState() error {
	path, err := pathForWrite(Config, "*")
	if err != nil {
		log.Errorf("configs cannot be written: %v", err)
	} else {
		log.Infof("configs will be written to %s", path)
	}
	path, err = pathForWrite(SavedGames, "*")
	if err != nil {
		log.Errorf("save games cannot be written: %v", err)
	} else {
		log.Infof("save games will be written to %s", path)
	}
	return nil
}

// ReadState loads the given state file and returns its contents.
func ReadState(kind StateKind, name string) ([]byte, error) {
	path, err := pathForRead(kind, name)
	if err != nil {
		// Remap to os.ErrNotExist so callers can deal with the error on their own.
		// This error is expected on first run, so it's just INFO.
		log.Infof("could not find path for folder%d/%s: %v", kind, name, err)
		return nil, os.ErrNotExist
	}
	return ioutil.ReadFile(path)
}

// MoveAwayState renames a detected-to-be-broken state file so it will not be used again.
func MoveAwayState(kind StateKind, name string) error {
	suffix := time.Now().UTC().Format(".2006-01-02T15-04-05Z")
	oldName, err := pathForRead(kind, name)
	if err != nil {
		return err
	}
	newName := oldName + suffix
	log.Errorf("renaming broken state file %s -> %v", oldName, newName)
	err = os.Rename(oldName, newName)
	if err == os.ErrNotExist {
		// Source file not therre? I guess that is fine too.
		return nil
	}
	return err
}

// writeState writes the given state file.
func writeState(kind StateKind, name string, data []byte) error {
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
