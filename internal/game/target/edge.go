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
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// EdgeTarget cleans up SetState events to match edges only.
type EdgeTarget struct {
	World  *engine.World
	Entity *engine.Entity

	Target mixins.TargetSelection

	State bool
}

func (t *EdgeTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	t.World = w
	t.Entity = e
	t.Target = mixins.ParseTarget(sp.Properties["target"])
	t.State = sp.Properties["initial_state"] == "true"
	mixins.SetStateOfTarget(t.World, t.Entity, t.Target, t.State)
	return nil
}

func (t *EdgeTarget) Despawn() {}

func (t *EdgeTarget) Update() {}

func (t *EdgeTarget) Touch(other *engine.Entity) {}

func (t *EdgeTarget) SetState(state bool) {
	if state == t.State {
		return
	}
	t.State = state
	mixins.SetStateOfTarget(t.World, t.Entity, t.Target, t.State)
}

func (t *EdgeTarget) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&EdgeTarget{})
}
