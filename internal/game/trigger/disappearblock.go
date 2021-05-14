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
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/constants"
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/level"
)

// DisappearBlock is an entity that disappears when touched and never reappears.
type DisappearBlock struct {
	World  *engine.World
	Entity *engine.Entity

	Disappearing bool
	AnimFrame    int
}

const (
	DisappearFrames         = 32
	DisappearSolidThreshold = 1
)

func (a *DisappearBlock) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	a.World = w
	a.Entity = e

	var err error
	e.Image, err = image.Load("sprites", "disappearblock.png")
	if err != nil {
		return err
	}
	e.Alpha = 0.0
	w.SetZIndex(e, constants.DisappearBlockZ)
	a.AnimFrame = DisappearFrames

	return nil
}

func (a *DisappearBlock) Despawn() {}

func (a *DisappearBlock) Update() {
	if !a.Disappearing {
		delta := a.Entity.Rect.Delta(a.World.Player.Rect)
		if delta.Norm1() <= 1 {
			// Touching the block from a face. Touching from a corner does not count.
			a.Disappearing = true
		}
	}
	if a.Disappearing {
		if a.AnimFrame > 0 {
			a.AnimFrame--
		}
	}
	a.Entity.Alpha = float64(a.AnimFrame) / DisappearFrames
	a.World.MutateContentsBool(a.Entity, level.PlayerSolidContents, a.AnimFrame >= DisappearSolidThreshold)
}

func (a *DisappearBlock) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&DisappearBlock{})
}
