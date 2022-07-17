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
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/palette"
)

func main() {
	log.Debugf("parsing flags...")
	flag.Parse(flag.NoConfig)
	log.Debugf("dumping palette LUTs...")
	err := palette.SaveCachedLUTs(640, 360, "assets/generated")
	if err != nil {
		log.Fatalf("failed dumping palettes: %v", err)
	}
	log.Debugf("done.")
}
