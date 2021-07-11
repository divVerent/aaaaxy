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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// QuestionBlock is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type QuestionBlock struct {
	World           *engine.World
	Entity          *engine.Entity
	PersistentState map[string]string

	Kaizo        bool
	Used         bool
	UsedImage    *ebiten.Image
	UseAnimFrame int

	Sound *sound.Sound
}

const (
	UseFramesPerPixel = 2
	UsePixels         = 4
)

func (q *QuestionBlock) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	q.World = w
	q.Entity = e
	q.PersistentState = s.PersistentState

	var err error
	w.SetSolid(e, true)
	w.SetOpaque(e, false)        // These shadows are annoying.
	e.Orientation = m.Identity() // Always show upright.
	q.Kaizo = s.Properties["kaizo"] == "true"
	q.Used = q.PersistentState["used"] == "true"
	q.UsedImage, err = image.Load("sprites", "exclamationblock.png")
	if err != nil {
		return err
	}
	if q.Used {
		e.Image = q.UsedImage
		q.UseAnimFrame = 2 * UseFramesPerPixel * UsePixels
	} else {
		if !q.Kaizo {
			e.Image, err = image.Load("sprites", "questionblock.png")
			if err != nil {
				return err
			}
		}
	}
	q.Sound, err = sound.Load("questionblock.ogg")
	if err != nil {
		return fmt.Errorf("could not load questionblock sound: %v", err)
	}
	return nil
}

func (q *QuestionBlock) Despawn() {}

func (q *QuestionBlock) isAbove(other *engine.Entity) bool {
	onGroundVec := m.Delta{DX: 0, DY: 1}
	if phys, ok := other.Impl.(interfaces.Physics); ok {
		onGroundVec = phys.ReadOnGroundVec()
	}
	return q.Entity.Rect.Delta(other.Rect).Dot(onGroundVec) < 0
}

func (q *QuestionBlock) Update() {
	if q.Used {
		if q.UseAnimFrame < UseFramesPerPixel*UsePixels {
			q.UseAnimFrame++
			if q.UseAnimFrame%UseFramesPerPixel == 0 {
				q.Entity.Rect.Origin.Y--
			}
		} else if q.UseAnimFrame < 2*UseFramesPerPixel*UsePixels {
			q.UseAnimFrame++
			if q.UseAnimFrame%UseFramesPerPixel == 0 {
				q.Entity.Rect.Origin.Y++
			}
		}
		return
	}
	if !q.Kaizo {
		return
	}
	q.World.SetSolid(q.Entity, q.isAbove(q.World.Player))
}

func (q *QuestionBlock) Touch(other *engine.Entity) {
	if other != q.World.Player {
		return
	}
	if q.Used {
		return
	}
	if !q.isAbove(other) {
		return
	}
	q.Used = true
	q.PersistentState["used"] = "true"
	q.Entity.Image = q.UsedImage
	q.UsedImage = nil
	q.World.SetSolid(q.Entity, true)
	q.Sound.Play()
}

func init() {
	engine.RegisterEntityType(&QuestionBlock{})
}
