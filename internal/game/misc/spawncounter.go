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

package misc

import (
	"fmt"
	"strconv"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
)

// SpawnCounter triggers a target if it's been spawned a certain amount of times.
type SpawnCounter struct{}

func (s *SpawnCounter) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	target := mixins.ParseTarget(sp.Properties["target"])
	state := sp.Properties["state"] != "false"
	divisorStr := sp.Properties["divisor"]
	divisor, err := strconv.Atoi(divisorStr)
	if err != nil {
		return fmt.Errorf("could not decode divisor %q: %v", divisorStr, err)
	}
	modulusStr := sp.Properties["modulus"]
	modulus, err := strconv.Atoi(modulusStr)
	if err != nil {
		return fmt.Errorf("could not decode modulus %q: %v", modulusStr, err)
	}
	count := 0
	countStr := sp.PersistentState["count"]
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("could not decode count %q: %v", countStr, err)
		}
	}
	count++
	sp.PersistentState["count"] = fmt.Sprint(count)
	if count%divisor == modulus {
		mixins.SetStateOfTarget(w, e, e, target, state)
	}
	return nil
}

func (s *SpawnCounter) Despawn() {}

func (s *SpawnCounter) Update() {}

func (s *SpawnCounter) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&SpawnCounter{})
}
