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

package trigger

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// SetState overrides the boolean state of a warpzone or entity.
type SetState struct {
	mixins.NonSolidTouchable
	World *engine.World

	Target string
	State  bool
}

func (s *SetState) Spawn(w *engine.World, sp *engine.Spawnable, e *engine.Entity) error {
	s.NonSolidTouchable.Init(w, e)
	s.World = w
	s.Target = sp.Properties["target"]
	s.State = sp.Properties["state"] == "true"
	if sp.Properties["initial_state"] != "" {
		s.apply(sp.Properties["initial_state"] == "true")
	}
	return nil
}

func (s *SetState) Despawn() {}

type stateSetter interface {
	SetState(state bool)
}

func (s *SetState) apply(state bool) {
	s.World.SetWarpZoneState(s.Target, state)
	for _, ent := range s.World.Entities {
		if ent.Name != s.Target {
			continue
		}
		setter, ok := ent.Impl.(stateSetter)
		if !ok {
			log.Panicf("Tried to set state of a non-supporting entity: %T, name: %v", s.Target)
		}
		setter.SetState(state)
	}
}

func (s *SetState) Touch(other *engine.Entity) {
	if other == s.World.Player {
		s.apply(s.State)
	}
}

func (s *SetState) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&SetState{})
}
