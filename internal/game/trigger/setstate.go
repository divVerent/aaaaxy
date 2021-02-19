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
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// SetState overrides the boolean state of a warpzone or entity.
type SetState struct {
	mixins.NonSolidTouchable
	World  *engine.World
	Entity *engine.Entity

	Target string
	State  bool
}

func (s *SetState) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.NonSolidTouchable.Init(w, e)
	s.World = w
	s.Entity = e
	s.Target = sp.Properties["target"]
	s.State = sp.Properties["state"] != "false"
	if sp.Properties["initial_state"] != "" {
		mixins.SetStateOfTarget(w, e, s.Target, sp.Properties["initial_state"] != "false") // Default true.
	}
	return nil
}

func (s *SetState) Despawn() {}

func (s *SetState) Touch(other *engine.Entity) {
	if other == s.World.Player {
		mixins.SetStateOfTarget(s.World, s.Entity, s.Target, s.State)
	}
}

func (s *SetState) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&SetState{})
}
