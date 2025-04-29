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

	"github.com/divVerent/aaaaxy/internal/animation"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/m"
)

// OneWay is an entity that can only be passed in one direction.
// It is implemented simply as being solid whenever it is on the left (or rotated direction) of the player.
type OneWay struct {
	World  *engine.World
	Entity *engine.Entity

	AllowedDirection m.Delta

	Anim animation.State
}

func (o *OneWay) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	o.World = w
	o.Entity = e

	e.ResizeImage = true
	w.SetOpaque(e, false)
	w.SetZIndex(e, constants.OneWayZ)

	o.AllowedDirection = e.Orientation.Apply(m.East())

	err := o.Anim.Init("oneway", map[string]*animation.Group{
		"idle": {
			Frames:        2,
			FrameInterval: 16,
			NextInterval:  32,
			NextAnim:      "idle",
		},
	}, "idle")
	if err != nil {
		return fmt.Errorf("could not initialize oneway animation: %w", err)
	}

	return nil
}

func (o *OneWay) Despawn() {}

func (o *OneWay) Update() {
	// I am solid if plP0..plP1 is entirely > myP0..myP1.
	myP0 := o.Entity.Rect.Origin.Delta(m.Pos{}).Dot(o.AllowedDirection)
	myP1 := o.Entity.Rect.OppositeCorner().Delta(m.Pos{}).Dot(o.AllowedDirection)
	plP0 := o.World.Player.Rect.Origin.Delta(m.Pos{}).Dot(o.AllowedDirection)
	plP1 := o.World.Player.Rect.OppositeCorner().Delta(m.Pos{}).Dot(o.AllowedDirection)
	o.World.MutateContentsBool(o.Entity, level.PlayerSolidContents, plP0 > myP0 && plP1 > myP0 && plP0 > myP1 && plP1 > myP1)

	// Animate.
	o.Anim.Update(o.Entity)
}

func (o *OneWay) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&OneWay{})
}
