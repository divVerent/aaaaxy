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

package main

import (
	"encoding/json"
	"os"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

func main() {
	log.Debugf("initializing VFS...")
	err := vfs.Init()
	if err != nil {
		log.Fatalf("could not initialize VFS: %v", err)
	}
	log.Debugf("parsing flags...")
	flag.Parse(flag.NoConfig)
	log.Debugf("loading level...")
	lvl, err := level.NewLoader("level").SkipComparingCheckpointLocations(true).Load()
	if err != nil {
		log.Fatalf("could not load level: %v", err)
	}
	j := json.NewEncoder(os.Stdout)
	j.SetIndent("", "\t")
	err = j.Encode(lvl.CheckpointLocations)
	if err != nil {
		log.Fatalf("could not dump locations: %v", err)
	}
	log.Debugf("done.")
}
