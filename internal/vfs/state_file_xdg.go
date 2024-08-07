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

//go:build !wasm && !windows && !android && !ios && !darwin
// +build !wasm,!windows,!android,!ios,!darwin

package vfs

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
)

func pathForReadRaw(kind StateKind, name string) ([]string, error) {
	switch kind {
	case Config:
		path, err := xdg.SearchConfigFile(filepath.Join(gameName, name))
		return []string{path}, err
	case SavedGames:
		path, err := xdg.SearchDataFile(filepath.Join(gameName, name))
		return []string{path}, err
	default:
		return nil, fmt.Errorf("searched for unsupported state kind: %d", kind)
	}
}

func pathForWriteRaw(kind StateKind, name string) (string, error) {
	switch kind {
	case Config:
		return xdg.ConfigFile(filepath.Join(gameName, name))
	case SavedGames:
		return xdg.DataFile(filepath.Join(gameName, name))
	default:
		return "", fmt.Errorf("searched for unsupported state kind: %d", kind)
	}
}
