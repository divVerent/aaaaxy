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

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// SpawnCounter triggers a target if it's been spawned a certain amount of times.
type SpawnCounter struct{}

func (s *SpawnCounter) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	var parseErr error
	state := propmap.ValueOrP(sp.Properties, "state", true, &parseErr)

	count := propmap.ValueOrP(sp.PersistentState, "count", 0, &parseErr)
	count++
	propmap.Set(sp.PersistentState, "count", count)

	for i := 1; ; i++ {
		suffix := ""
		if i > 1 {
			suffix = fmt.Sprint(i)
		}
		divisor := propmap.ValueOrP(sp.Properties, "divisor"+suffix, -1, &parseErr)
		if divisor < 0 {
			break
		}
		modulus := propmap.ValueP(sp.Properties, "modulus"+suffix, 0, &parseErr)
		target := mixins.ParseTarget(propmap.ValueP(sp.Properties, "target"+suffix, "", &parseErr))
		if count%divisor == modulus {
			mixins.SetStateOfTarget(w, e, e, target, state)
		}
	}

	return parseErr
}

func (s *SpawnCounter) Despawn() {}

func (s *SpawnCounter) Update() {}

func (s *SpawnCounter) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&SpawnCounter{})
}
