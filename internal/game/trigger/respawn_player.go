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
	"fmt"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/picture"
)

// RespawnPlayer respawns the player when touched.
type RespawnPlayer struct {
	World *engine.World
}

// Let's do a somewhat forgiving hitbox.
const (
	RespawnHitboxRemoveAtSides = 4 // Leaves 8px between two spikes, player doesn't fit there.
	RespawnHitboxRemoveAtTop   = 4 // Only keep the bottom 12px.
)

func (r *RespawnPlayer) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	r.World = w
	var err error
	e.Image, err = picture.Load("sprites", "spike.png")
	if err != nil {
		return fmt.Errorf("could not load spike image: %r", err)
	}
	shrinkerTopLeft := e.Orientation.Apply(m.Delta{DX: -RespawnHitboxRemoveAtTop, DY: RespawnHitboxRemoveAtSides})
	shrinkerBottomRight := e.Orientation.Apply(m.Delta{DX: 0, DY: -RespawnHitboxRemoveAtSides})
	org := e.Rect.Origin
	e.Rect = e.Rect.ShrinkInDirection(shrinkerTopLeft).ShrinkInDirection(shrinkerBottomRight)
	e.RenderOffset = org.Delta(e.Rect.Origin)
	w.SetSolid(e, true)
	w.MutateContentsBool(e, level.PlayerSteppableSolidContents, false)
	w.SetZIndex(e, constants.RespawnPlayerZ)
	return nil
}

func (r *RespawnPlayer) Despawn() {}

func (r *RespawnPlayer) Update() {}

func (r *RespawnPlayer) Touch(other *engine.Entity) {
	if other != r.World.Player {
		return
	}
	r.World.ScheduleRespawnPlayer(r.World.PlayerState.LastCheckpoint(), false)
}

func init() {
	engine.RegisterEntityType(&RespawnPlayer{})
}
