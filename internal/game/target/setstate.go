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
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// SetStateTarget overrides the boolean state of a warpzone or entity.
type SetStateTarget struct {
	World  *engine.World
	Entity *engine.Entity

	Target           string
	ExclusiveTargets []string
	State            bool
}

func (s *SetStateTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.World = w
	s.Entity = e
	s.Target = sp.Properties["target"]
	s.State = sp.Properties["state"] != "false"
	s.ExclusiveTargets = strings.Split(sp.Properties["exclusive_targets"], " ")
	if sp.Properties["initial_state"] != "" {
		mixins.SetStateOfTarget(w, e, s.Target, sp.Properties["initial_state"] != "false") // Default true.
	}
	return nil
}

func (s *SetStateTarget) Despawn() {}

func (s *SetStateTarget) Update() {}

func (s *SetStateTarget) SetState(state bool) {
	// Turn my targets to s.State if state, to !s.State if !state.
	mixins.SetStateOfTarget(s.World, s.Entity, s.Target, s.State == state)

	// Turn off all my "siblings". Can be used to make mutually exclusive targets.
	for _, t := range s.ExclusiveTargets {
		mixins.SetStateOfTarget(s.World, s.Entity, t, false)
	}
}

func (s *SetStateTarget) Touch(other *engine.Entity) {}

func (s *SetStateTarget) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&SetStateTarget{})
}
