// Copyright 2022 Google LLC
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

//go:build android
// +build android

package vfs

import (
	"fmt"
	"path/filepath"

	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	filesDir string
)

func SetFilesDir(dir string) {
	filesDir = dir
}

func pathForReadRaw(kind StateKind, name string) (string, error) {
	return pathForWrite(kind, name)
}

func pathForWriteRaw(kind StateKind, name string) (string, error) {
	if filesDir == "" {
		log.Fatalf("tried to access data but SetFilesDir was not called yet")
	}
	switch kind {
	case Config:
		return filepath.Join(filesDir, "config", name), nil
	case SavedGames:
		return filepath.Join(filesDir, "save", name), nil
	default:
		return "", fmt.Errorf("searched for unsupported state kind: %d", kind)
	}
}
