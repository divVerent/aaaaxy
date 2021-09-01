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
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
)

// RespawnPlayer respawns the player when touched.
type RespawnPlayer struct {
	World *engine.World
}

// Let's do a somewhat forgiving hitbox.
const (
	RespawnHitboxBorder = 4 // Actually only the 8x8 center square of the 16x16 sprite.
)

func (r *RespawnPlayer) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	r.World = w
	var err error
	e.Image, err = image.Load("sprites", "spike.png")
	if err != nil {
		return fmt.Errorf("could not load spike image: %r", err)
	}
	e.RenderOffset = m.Delta{DX: -RespawnHitboxBorder, DY: -RespawnHitboxBorder}
	e.Rect.Origin = e.Rect.Origin.Sub(e.RenderOffset)
	e.Rect.Size = e.Rect.Size.Add(e.RenderOffset.Mul(2))
	w.SetSolid(e, true)
	w.SetZIndex(e, constants.RespawnPlayerZ)
	return nil
}

func (r *RespawnPlayer) Despawn() {}

func (r *RespawnPlayer) Update() {}

func (r *RespawnPlayer) Touch(other *engine.Entity) {
	if other != r.World.Player {
		return
	}
	r.World.RespawnPlayer(r.World.PlayerState.LastCheckpoint())
}

func init() {
	engine.RegisterEntityType(&RespawnPlayer{})
}
