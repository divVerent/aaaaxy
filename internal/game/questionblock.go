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

package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/image"
	m "github.com/divVerent/aaaaaa/internal/math"
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
}

const (
	UseFramesPerPixel = 2
	UsePixels         = 4
)

func (q *QuestionBlock) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	q.World = w
	q.Entity = e
	q.PersistentState = s.PersistentState

	var err error
	e.Solid = true
	e.Opaque = false             // These shadows are annoying.
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
	return nil
}

func (q *QuestionBlock) Despawn() {}

func (q *QuestionBlock) isAbove(other *engine.Entity) bool {
	return q.Entity.Rect.OppositeCorner().Y < other.Rect.Origin.Y
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
	q.Entity.Solid = q.isAbove(q.World.Player)
}

func (q *QuestionBlock) Touch(other *engine.Entity) {
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
	q.Entity.Solid = true
	// TODO animate up and down?
}

func (q *QuestionBlock) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&QuestionBlock{})
}
