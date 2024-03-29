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

package target

import (
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// SetStateTarget overrides the boolean state of a warpzone or entity.
type SetStateTarget struct {
	World  *engine.World
	Entity *engine.Entity

	Target mixins.TargetSelection
	State  bool
}

func (s *SetStateTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	s.World = w
	s.Entity = e
	var parseErr error
	s.Target = mixins.ParseTarget(propmap.ValueP(sp.Properties, "target", "", &parseErr))
	s.State = propmap.ValueOrP(sp.Properties, "state", true, &parseErr)
	return parseErr
}

func (s *SetStateTarget) Despawn() {}

func (s *SetStateTarget) Update() {}

func (s *SetStateTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	// Turn my targets to s.State if state, to !s.State if !state.
	mixins.SetStateOfTarget(s.World, originator, s.Entity, s.Target, s.State == state)
}

func (s *SetStateTarget) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&SetStateTarget{})
}
