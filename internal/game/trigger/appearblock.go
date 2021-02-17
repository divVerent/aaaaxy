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
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// AppearBlock is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type AppearBlock struct {
	World  *engine.World
	Entity *engine.Entity

	AnimFrame int
}

const (
	AppearFrames         = 16
	AppearXDistance      = 2 * level.TileSize
	AppearYDistance      = level.TileSize / 4
	AppearSolidThreshold = 12
)

func (a *AppearBlock) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	a.World = w
	a.Entity = e

	var err error
	e.Image, err = image.Load("sprites", "appearblock.png")
	if err != nil {
		return err
	}
	e.Opaque = false
	e.Solid = false
	e.Alpha = 0.0
	e.ZIndex = 1 // Above normal sprites.

	return nil
}

func (a *AppearBlock) Despawn() {}

func (a *AppearBlock) Update() {
	delta := a.Entity.Rect.Delta(a.World.Player.Rect)
	if delta.DY > 0 && delta.DX <= AppearXDistance && delta.DX >= -AppearXDistance && delta.DY <= AppearYDistance && delta.DY >= -AppearYDistance {
		if a.AnimFrame < AppearFrames {
			a.AnimFrame++
		}
	} else {
		if a.AnimFrame > 0 {
			a.AnimFrame--
		}
	}
	a.Entity.Alpha = float64(a.AnimFrame) / AppearFrames
	// Make nonsolid if inside (to unstick player).
	a.Entity.Solid = a.AnimFrame >= AppearSolidThreshold && delta.DY > 0
}

func (a *AppearBlock) Touch(other *engine.Entity) {}

func (a *AppearBlock) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&AppearBlock{})
}
