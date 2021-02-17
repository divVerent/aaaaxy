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

	"github.com/divVerent/aaaaaa/internal/animation"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// OneWay is an entity that can only be passed in one direction.
// It is implemented simply as being solid whenever it is on the left (or rotated direction) of the player.
type OneWay struct {
	World  *engine.World
	Entity *engine.Entity

	AllowedDirection m.Delta

	Anim animation.State
}

func (o *OneWay) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	o.World = w
	o.Entity = e

	e.ResizeImage = true
	e.Opaque = false

	o.AllowedDirection = e.Orientation.Apply(m.East())

	o.Anim.Init("oneway", map[string]*animation.Group{
		"idle": {
			Frames:        2,
			FrameInterval: 16,
			NextInterval:  32,
			NextAnim:      "idle",
		},
	}, "idle")

	return nil
}

func (o *OneWay) Despawn() {}

func (o *OneWay) Update() {
	// I am solid if plP0..plP1 is entirely > myP0..myP1.
	myP0 := o.Entity.Rect.Origin.Delta(m.Pos{}).Dot(o.AllowedDirection)
	myP1 := o.Entity.Rect.OppositeCorner().Delta(m.Pos{}).Dot(o.AllowedDirection)
	plP0 := o.World.Player.Rect.Origin.Delta(m.Pos{}).Dot(o.AllowedDirection)
	plP1 := o.World.Player.Rect.OppositeCorner().Delta(m.Pos{}).Dot(o.AllowedDirection)
	o.Entity.Solid = plP0 > myP0 && plP1 > myP0 && plP0 > myP1 && plP1 > myP1

	// Animate.
	o.Anim.Update(o.Entity)
}

func (o *OneWay) Touch(other *engine.Entity) {}

func (o *OneWay) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&OneWay{})
}
