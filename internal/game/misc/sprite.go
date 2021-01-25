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

package misc

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/image"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// Sprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
// Can be toggled from outside.
type Sprite struct {
	Entity *engine.Entity

	Image  *ebiten.Image
	Solid  bool
	Opaque bool
}

func (s *Sprite) Spawn(w *engine.World, sp *engine.Spawnable, e *engine.Entity) error {
	s.Entity = e

	var err error
	directory := sp.Properties["image_dir"]
	if directory == "" {
		directory = "sprites"
	}
	s.Image, err = image.Load(directory, sp.Properties["image"])
	if err != nil {
		return err
	}
	e.ResizeImage = true
	s.Solid = sp.Properties["solid"] == "true"
	s.Opaque = sp.Properties["opaque"] == "true"
	if sp.Properties["alpha"] != "" {
		e.Alpha, err = strconv.ParseFloat(sp.Properties["alpha"], 64)
		if err != nil {
			return fmt.Errorf("could not decode alpha %q: %v", sp.Properties["alpha"], err)
		}
	}
	if sp.Properties["z_index"] != "" {
		e.ZIndex, err = strconv.Atoi(sp.Properties["z_index"])
		if err != nil {
			return fmt.Errorf("could not decode z index %q: %v", sp.Properties["z_index"], err)
		}
	}
	if sp.Properties["no_transform"] == "true" {
		// Undo transform of orientation by tile.
		e.Orientation = sp.Orientation
	}
	if e.Transform.Determinant() < 0 {
		// e.Orientation: in-editor transform. Applied first.
		// Normally the formula is e.Transform.Inverse().Concat(e.Orientation).
		// Add an FlipX() between the two to "undo" any sense difference in the editor.
		// This flips the view on the _level editor_ X axis.
		switch sp.Properties["no_flip"] {
		case "x":
			e.Orientation = e.Transform.Inverse().Concat(m.FlipX()).Concat(e.Orientation)
		case "y":
			e.Orientation = e.Transform.Inverse().Concat(m.FlipY()).Concat(e.Orientation)
		}
	}
	initialState := sp.Properties["initial_state"] != "false"
	s.SetState(initialState)
	return nil
}

func (s *Sprite) SetState(state bool) {
	if state {
		s.Entity.Image = s.Image
		s.Entity.Solid = s.Solid
		s.Entity.Opaque = s.Opaque
	} else {
		s.Entity.Image = nil
		s.Entity.Solid = false
		s.Entity.Opaque = false
	}
}

func (s *Sprite) Despawn() {}

func (s *Sprite) Update() {}

func (s *Sprite) Touch(other *engine.Entity) {}

func (s *Sprite) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&Sprite{})
}
