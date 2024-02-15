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

package vfs

import (
	"os"
	"path/filepath"

	"github.com/divVerent/aaaaxy/internal/log"
)

// exeDir is the directory of the current executable.
var exeDir string = ""

func initExeDir() {
	exePath, err := os.Executable()
	if err != nil {
		log.Warningf("could not find path to executable: %v; using current working directory instead", err)
		return
	}
	exeDir = filepath.Dir(exePath)
	log.Infof("found executable directory: %v", exeDir)
}
